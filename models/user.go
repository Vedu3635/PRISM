package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"size:255;uniqueIndex"`
	Username     string    `gorm:"size:50;uniqueIndex"`
	FullName     string    `gorm:"size:100"`
	PasswordHash string
	AvatarURL    *string
	Phone        *string
	CurrencyPref string `gorm:"default:INR"`

	IsVerified bool `gorm:"default:false"`
	IsDeleted  bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
