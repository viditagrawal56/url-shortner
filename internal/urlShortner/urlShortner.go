package urlShortner

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/viditagrawal56/url-shortner/internal/models"
	"gorm.io/gorm"
)

var (
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

// TODO: Improve the shortining algorithm
func (s *URLShortnerService) generateUniqueShortCode() (string, error) {
	for attempts := 0; attempts < 5; attempts++ {
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
