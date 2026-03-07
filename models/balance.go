package models

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID    uuid.UUID `gorm:"index"`
	FromUserID uuid.UUID `gorm:"index"`
	ToUserID   uuid.UUID `gorm:"index"`

	NetAmount float64
	Currency  string

	LastUpdated time.Time
}
