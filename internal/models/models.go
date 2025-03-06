package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()" json:"updated_at"`

	// Associations
	ShortURLs []ShortURL `gorm:"foreignKey:UserID" json:"short_urls,omitempty"`
}

type ShortURL struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	ShortCode      string     `gorm:"type:varchar(7);not null;unique" json:"short_code"`
	OriginalURL    string     `gorm:"type:text;not null" json:"original_url"`
	RequiresAuth   bool       `gorm:"type:boolean;not null;default:false" json:"requires_auth"`
	NotifyOnAccess bool       `gorm:"type:boolean;not null;default:false" json:"notify_on_access"`
	Active         bool       `gorm:"type:boolean;not null;default:false" json:"active"`
	ValidFrom      *time.Time `gorm:"type:timestamp with time zone" json:"valid_from"`
	ValidUntil     *time.Time `gorm:"type:timestamp with time zone" json:"valid_until"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:timestamp with time zone;not null;default:now()" json:"updated_at"`

	// Associations
	User             User              `gorm:"foreignKey:UserID" json:"-"`
	AuthorizedEmails []AuthorizedEmail `gorm:"foreignKey:ShortURLID" json:"authorized_emails,omitempty"`
}

type AuthorizedEmail struct {
	ShortURLID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"short_url_id"`
	Email      string    `gorm:"type:varchar(255);primaryKey;not null" json:"email"`

	// Associations
	ShortURL ShortURL `gorm:"foreignKey:ShortURLID" json:"-"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
