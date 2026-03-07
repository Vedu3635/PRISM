package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;index"`

	Type string

	Title string
	Body  string

	EntityType *string
	EntityID   *uuid.UUID

	IsRead bool `gorm:"default:false"`

	CreatedAt time.Time
}
