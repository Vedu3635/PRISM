package models

import (
	"time"

	"github.com/google/uuid"
)

type SettlementSplit struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	SettlementID uuid.UUID `gorm:"type:uuid;index"`
	SplitID      uuid.UUID `gorm:"type:uuid;index"`

	AppliedAmount float64 `gorm:"type:numeric(12,2)"`

	Cleared bool `gorm:"default:false"`

	CreatedAt time.Time
}
