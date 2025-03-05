package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
