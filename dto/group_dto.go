package dto

import "github.com/google/uuid"

type CreateGroupRequest struct {
	CreatedBy   uuid.UUID `json:"created_by"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Type        string    `json:"type"`
	Currency    string    `json:"currency"`
}

type AddGroupMemberRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}
