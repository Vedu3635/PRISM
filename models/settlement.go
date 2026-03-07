package models

import (
	"time"

	"github.com/google/uuid"
)

type Settlement struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID    uuid.UUID `gorm:"index"`
	FromUserID uuid.UUID `gorm:"index"`
	ToUserID   uuid.UUID `gorm:"index"`

	Amount   float64
	Currency string

	Status string `gorm:"default:pending"`

	PaymentMethod *string
	PaymentID     *uuid.UUID

	Note      *string
	SettledAt *time.Time
	CreatedAt time.Time
}
