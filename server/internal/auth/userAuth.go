package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/viditagrawal56/url-shortner/internal/config"
	"github.com/viditagrawal56/url-shortner/internal/db"
	"github.com/viditagrawal56/url-shortner/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrMissingToken       = errors.New("missing authorization token")
	ErrInvalidToken       = errors.New("invalid token format")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidSignature   = errors.New("invalid token signature")
)

type Service struct {
	db  *db.Database
	cfg *config.Config
}

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	UserIDKey contextKey = "userID"
	EmailKey  contextKey = "email"
)

const MagicLinkTokenExpiration = 15 * time.Minute

func New(database *db.Database, cfg *config.Config) *Service {
	return &Service{
		db:  database,
		cfg: cfg,
	}
}

func (s *Service) generateJWTToken(userID uuid.UUID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.Auth.TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.Auth.JWTSecret))
}

func (s *Service) RegisterUser(email, password string) error {
	if email == "" || password == "" {
		return ErrInvalidInput
	}

	// Check if user already exists
	if result := s.db.GetDB().Where("email = ?", email).First(&models.User{}); result.Error == nil {
		return ErrUserExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.db.GetDB().Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *Service) LoginUser(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", ErrInvalidInput
	}

	// Find user
	var user models.User
	result := s.db.GetDB().Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("database error: %w", result.Error)
	}

	// Verify Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.generateJWTToken(user.ID, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return token, nil
}

func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeAuthError(w, ErrMissingToken)
			return
		}

		// Check if the Authorization header is in the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeAuthError(w, ErrInvalidToken)
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(s.cfg.Auth.JWTSecret), nil
		})

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrSignatureInvalid):
				writeAuthError(w, ErrInvalidSignature)
			case errors.Is(err, jwt.ErrTokenExpired):
				writeAuthError(w, ErrExpiredToken)
			default:
				writeAuthError(w, ErrInvalidToken)
			}
			return
		}

		if !token.Valid {
			writeAuthError(w, ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)

		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) GenerateMagicLinkToken(email, shortCode string) (string, error) {
	if email == "" || shortCode == "" {
		return "", ErrInvalidInput
	}

	tokenBytes := make([]byte, 32)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	tokenID := uuid.NewString()

	claims := jwt.MapClaims{
		"email":      email,
		"short_code": shortCode,
		"token_id":   tokenID,
		"exp":        time.Now().Add(MagicLinkTokenExpiration).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.cfg.Auth.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	// Store the token in database
	temporaryToken := models.TemporaryToken{
		Email:     email,
		ShortCode: shortCode,
		Token:     tokenID,
		ExpiresAt: time.Now().Add(MagicLinkTokenExpiration),
	}

	if err := s.db.GetDB().Create(&temporaryToken).Error; err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	return tokenString, nil
}

func (s *Service) VerifyMagicLinkToken(tokenString string) (string, string, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return "", "", ErrInvalidSignature
		default:
			return "", "", ErrInvalidToken
		}
	}

	if !token.Valid {
		return "", "", ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", ErrInvalidToken
	}

	shortCode, ok := claims["short_code"].(string)
	if !ok {
		return "", "", ErrInvalidToken
	}

	tokenID, ok := claims["token_id"].(string)
	if !ok {
		return "", "", ErrInvalidToken
	}

	var temporaryToken models.TemporaryToken
	result := s.db.GetDB().Where("token = ? AND email = ? AND short_code = ? AND used_at IS NULL AND expires_at > ?",
		tokenID, email, shortCode, time.Now()).First(&temporaryToken)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", "", ErrInvalidToken
	}

	if result.Error != nil {
		return "", "", fmt.Errorf("failed to verify token: %w", result.Error)
	}

	now := time.Now()
	temporaryToken.UsedAt = &now
	if err := s.db.GetDB().Save(&temporaryToken).Error; err != nil {
		return "", "", fmt.Errorf("failed to mark token as used: %w", err)
	}

	return email, shortCode, nil
}

func writeAuthError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err {
	case ErrMissingToken:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authorization header required"})
	case ErrInvalidToken:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token format"})
	case ErrExpiredToken:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token has expired"})
	case ErrInvalidSignature:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token signature"})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
	}
}
