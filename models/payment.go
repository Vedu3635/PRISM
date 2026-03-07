package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	SettlementID uuid.UUID `gorm:"type:uuid;index"`

	Gateway        string
	GatewayTxnID   string `gorm:"uniqueIndex"`
	GatewayOrderID *string

	Amount   float64 `gorm:"type:numeric(12,2)"`
	Currency string

	Status string

	GatewayResponse map[string]interface{} `gorm:"type:jsonb"`

	IdempotencyKey string `gorm:"uniqueIndex"`

	CreatedAt time.Time
}
