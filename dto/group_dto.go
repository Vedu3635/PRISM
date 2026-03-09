package dto

import "github.com/google/uuid"

type CreateGroupRequest struct {
	CreatedBy   uuid.UUID `json:"created_by"   binding:"required"`
	Name        string    `json:"name"         binding:"required"`
	Description *string   `json:"description"`
	Type        string    `json:"type"         binding:"required"`
	Currency    string    `json:"currency"     binding:"required"`
}

type UpdateGroupRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	Currency    *string `json:"currency"`
}

type AddGroupMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Role   string    `json:"role"    binding:"required,oneof=admin member"`
}

type LeaveGroupRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
