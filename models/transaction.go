package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID uuid.UUID `gorm:"type:uuid;index"`
	PaidBy  uuid.UUID `gorm:"type:uuid;index"`

	Title    string
	Amount   float64
	Currency string

	Category  *string
	SplitType string

	Notes      *string
	ReceiptURL *string

	Status string `gorm:"default:active"`

	Metadata map[string]interface{} `gorm:"type:jsonb"`

	TransactedAt time.Time
	CreatedAt    time.Time
}
