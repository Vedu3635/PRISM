package dto

import "github.com/google/uuid"

type ParticipantSplit struct {
	UserID     uuid.UUID `json:"user_id"`
	OwedAmount float64   `json:"owed_amount"`
}

type CreateTransactionRequest struct {
	GroupID      uuid.UUID          `json:"group_id"`
	PaidBy       uuid.UUID          `json:"paid_by"`
	Title        string             `json:"title"`
	Amount       float64            `json:"amount"`
	SplitType    string             `json:"split_type"`
	Participants []ParticipantSplit `json:"participants"`
}
