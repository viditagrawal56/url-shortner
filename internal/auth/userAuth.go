package auth

import (
	"context"
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
	ErrInvalidSignature   = errors.New("invlaid token signature")
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
	userIDKey contextKey = "userID"
	emailKey  contextKey = "email"
)

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.Auth.JWTSecret))
}

func (s *Service) RegisterUser(email, password string) error {
	if email == "" || password == "" {
		return ErrInvalidInput
	}

	// Check if user already exists
	var existingUser models.User
	if result := s.db.GetDB().Where("email = ?", email).First(&existingUser); result.Error == nil {
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
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", ErrInvalidCredentials
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
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, emailKey, claims.Email)

		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value("userID").(uint)
	return id, ok
}
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("email").(string)
	return email, ok
}

func writeAuthError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err {
	case ErrMissingToken:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authorization header required"})
	case ErrInvalidToken:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invlaid token format"})
	case ErrExpiredToken:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token has expired"})
	case ErrInvalidSignature:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invliad token signature"})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
	}
}
