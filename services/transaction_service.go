package services

import (
	"github.com/Vedu3635/PRISM.git/database"
	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/models"
	"github.com/google/uuid"
)

func CreateTransaction(req dto.CreateTransactionRequest) (*models.Transaction, error) {

	db := database.DB

	tx := models.Transaction{
		ID:        uuid.New(),
		GroupID:   req.GroupID,
		PaidBy:    req.PaidBy,
		Title:     req.Title,
		Amount:    req.Amount,
		SplitType: req.SplitType,
	}

	if err := db.Create(&tx).Error; err != nil {
		return nil, err
	}

	for _, p := range req.Participants {

		split := models.TransactionSplit{
			ID:            uuid.New(),
			TransactionID: tx.ID,
			UserID:        p.UserID,
			OwedAmount:    p.OwedAmount,
		}

		if err := db.Create(&split).Error; err != nil {
			return nil, err
		}

		if p.UserID != req.PaidBy {
			updateBalance(req.GroupID, p.UserID, req.PaidBy, p.OwedAmount)
		}
	}

	return &tx, nil
}

func updateBalance(groupID, fromUser, toUser uuid.UUID, amount float64) {

	db := database.DB

	var balance models.Balance

	err := db.Where(
		"group_id = ? AND from_user_id = ? AND to_user_id = ?",
		groupID, fromUser, toUser,
	).First(&balance).Error

	if err != nil {

		balance = models.Balance{
			ID:         uuid.New(),
			GroupID:    groupID,
			FromUserID: fromUser,
			ToUserID:   toUser,
			NetAmount:  amount,
		}

		db.Create(&balance)
		return
	}

	balance.NetAmount += amount
	db.Save(&balance)
}
