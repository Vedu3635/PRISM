package dto

import "github.com/google/uuid"

type ParticipantSplit struct {
	UserID     uuid.UUID `json:"user_id"     binding:"required"`
	OwedAmount float64   `json:"owed_amount"  binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	GroupID      uuid.UUID          `json:"group_id"    binding:"required"`
	PaidBy       uuid.UUID          `json:"paid_by"     binding:"required"`
	Title        string             `json:"title"       binding:"required"`
	Amount       float64            `json:"amount"      binding:"required,gt=0"`
	Currency     string             `json:"currency"    binding:"required"`
	Category     *string            `json:"category"`
	SplitType    string             `json:"split_type"  binding:"required,oneof=equal exact percentage"`
	Notes        *string            `json:"notes"`
	ReceiptURL   *string            `json:"receipt_url"`
	Participants []ParticipantSplit `json:"participants" binding:"required,min=1,dive"`
}

type UpdateTransactionRequest struct {
	Title      *string `json:"title"`
	Category   *string `json:"category"`
	Notes      *string `json:"notes"`
	ReceiptURL *string `json:"receipt_url"`
}
