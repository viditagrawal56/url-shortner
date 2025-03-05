package auth

import (
	"errors"
	"fmt"

	"github.com/viditagrawal56/url-shortner/internal/config"
	"github.com/viditagrawal56/url-shortner/internal/db"
	"github.com/viditagrawal56/url-shortner/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidInput = errors.New("invalid input")
)

type Service struct {
	db  *db.Database
	cfg *config.Config
}

func New(database *db.Database, cfg *config.Config) *Service {
	return &Service{
		db:  database,
		cfg: cfg,
	}
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
