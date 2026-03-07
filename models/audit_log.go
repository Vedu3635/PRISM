package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	ActorID uuid.UUID `gorm:"type:uuid;index"`

	Action     string
	EntityType string
	EntityID   uuid.UUID `gorm:"index"`

	OldValue map[string]interface{} `gorm:"type:jsonb"`
	NewValue map[string]interface{} `gorm:"type:jsonb"`

	IPAddress *string
	RequestID *string

	CreatedAt time.Time
}
