package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedBy uuid.UUID `gorm:"type:uuid"`

	Name        string
	Description *string

	Type     string
	Currency string `gorm:"default:INR"`

	InviteCode string `gorm:"uniqueIndex"`

	IsActive   bool `gorm:"default:true"`
	IsPersonal bool `gorm:"default:false"`

	CreatedAt time.Time
}
