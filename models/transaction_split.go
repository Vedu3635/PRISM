package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionSplit struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	TransactionID uuid.UUID `gorm:"type:uuid;index"`
	UserID        uuid.UUID `gorm:"type:uuid;index"`

	OwedAmount float64
	Percentage *float64

	IsSettled bool `gorm:"default:false"`
	SettledAt *time.Time
}
