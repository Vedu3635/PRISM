package models

import (
	"time"

	"github.com/google/uuid"
)

type GroupMember struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID uuid.UUID `gorm:"type:uuid;index"`
	UserID  uuid.UUID `gorm:"type:uuid;index"`

	Role     string `gorm:"default:member"`
	Nickname *string

	JoinedAt time.Time
	LeftAt   *time.Time
}
