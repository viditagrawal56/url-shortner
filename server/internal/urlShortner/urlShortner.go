package urlShortner

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/viditagrawal56/url-shortner/internal/models"
	"gorm.io/gorm"
)

var (
	ErrURLNotFound            = errors.New("short URL not found")
	ErrURLExpired             = errors.New("URL has expired")
	ErrURLNotYetValid         = errors.New("URL is not yet valid")
	ErrAuthenticationRequired = errors.New("authentication required to access this URL")
	ErrUnauthorizedAccess     = errors.New("you are not authorized to access this URL")
)

const (
	charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodeLength = 7
)

type URLShortnerService struct {
	db *gorm.DB
}

func NewUrlShortnerService(db *gorm.DB) *URLShortnerService {
	return &URLShortnerService{
		db: db,
	}
}

func (s *URLShortnerService) CreateShortURL(userID uuid.UUID, originalURL string, options models.ShortURLOptions) (*models.ShortURL, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	shortCode, err := s.generateUniqueShortCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	shortURL := &models.ShortURL{
		UserID:         userID,
		ShortCode:      shortCode,
		OriginalURL:    originalURL,
		RequiresAuth:   options.RequiresAuth,
		NotifyOnAccess: options.NotifyOnAccess,
		Active:         true,
		ValidFrom:      options.ValidFrom,
		ValidUntil:     options.ValidUntil,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(shortURL).Error; err != nil {
			return fmt.Errorf("failed to create short URL: %w", err)
		}

		if options.RequiresAuth && len(options.AuthorizedEmails) > 0 {
			var authorizedEmails []models.AuthorizedEmail
			for _, email := range options.AuthorizedEmails {
				authorizedEmails = append(authorizedEmails, models.AuthorizedEmail{
					ShortURLID: shortURL.ID,
					Email:      email,
				})
			}

			if err := tx.Create(&authorizedEmails).Error; err != nil {
				return fmt.Errorf("failed to save authorized emails: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *URLShortnerService) ResolveShortURL(shortCode string, email string) (*models.ShortURL, error) {
	var shortURL models.ShortURL
	if err := s.db.Preload("AuthorizedEmails").
		Where("short_code = ? AND active = true", shortCode).
		First(&shortURL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to fetch short URL: %w", err)
	}

	// Check if URL is valid based on time constraints
	now := time.Now()
	if shortURL.ValidFrom != nil && now.Before(*shortURL.ValidFrom) {
		return nil, ErrURLNotYetValid
	}

	if shortURL.ValidUntil != nil && now.After(*shortURL.ValidUntil) {
		return nil, ErrURLExpired
	}

	if shortURL.RequiresAuth {
		if email == "" {
			return nil, ErrAuthenticationRequired
		}

		isAuthorized := false
		for _, authorizedEmail := range shortURL.AuthorizedEmails {
			if authorizedEmail.Email == email {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			return nil, ErrUnauthorizedAccess
		}
	}

	return &shortURL, nil
}

func (s *URLShortnerService) GetUserShortURLs(userID uuid.UUID) ([]models.ShortURL, error) {
	var shortURLs []models.ShortURL

	if err := s.db.Preload("AuthorizedEmails").Where("user_id = ?", userID).Find(&shortURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user's short URLs: %w", err)
	}

	return shortURLs, nil
}

// TODO: Improve the shortining algorithm
func (s *URLShortnerService) generateUniqueShortCode() (string, error) {
	for range 5 {
		shortCode, err := generateRandomString(shortCodeLength)
		if err != nil {
			return "", err
		}

		//Check if short code already exists
		var count int64
		if err := s.db.Model(&models.ShortURL{}).Where("short_code = ?", shortCode).Count(&count).Error; err != nil {
			return "", fmt.Errorf("failed to check short code uniqueness: %w", err)
		}

		if count == 0 {
			return shortCode, nil
		}
	}

	return "", errors.New("failed to generate a unique short code after multiple attempts")
}

func generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}
