package services

import (
	"errors"
	"time"

	"github.com/Vedu3635/PRISM.git/database"
	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTransaction(req dto.CreateTransactionRequest) (*models.Transaction, error) {
	db := database.DB

	// Validate split amounts sum to total
	var total float64
	for _, p := range req.Participants {
		total += p.OwedAmount
	}
	if total != req.Amount {
		return nil, errors.New("participant owed amounts must sum to total transaction amount")
	}

	var createdTx *models.Transaction

	// Wrap in DB transaction so splits + balances are atomic
	err := db.Transaction(func(tx *gorm.DB) error {
		transaction := models.Transaction{
			ID:           uuid.New(),
			GroupID:      req.GroupID,
			PaidBy:       req.PaidBy,
			Title:        req.Title,
			Amount:       req.Amount,
			Currency:     req.Currency,
			Category:     req.Category,
			SplitType:    req.SplitType,
			Notes:        req.Notes,
			ReceiptURL:   req.ReceiptURL,
			Status:       "active",
			TransactedAt: time.Now(),
			CreatedAt:    time.Now(),
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		for _, p := range req.Participants {
			split := models.TransactionSplit{
				ID:            uuid.New(),
				TransactionID: transaction.ID,
				UserID:        p.UserID,
				OwedAmount:    p.OwedAmount,
			}

			if err := tx.Create(&split).Error; err != nil {
				return err
			}

			if p.UserID != req.PaidBy {
				if err := updateBalance(tx, req.GroupID, p.UserID, req.PaidBy, p.OwedAmount); err != nil {
					return err
				}
			}
		}

		createdTx = &transaction
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTx, nil
}

func GetTransactions() ([]models.Transaction, error) {
	db := database.DB

	var transactions []models.Transaction
	if err := db.Where("status = ?", "active").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	db := database.DB

	var transaction models.Transaction
	if err := db.First(&transaction, "id = ? AND status = ?", id, "active").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}

func GetTransactionsByGroup(groupID uuid.UUID) ([]models.Transaction, error) {
	db := database.DB

	var transactions []models.Transaction
	if err := db.Where("group_id = ? AND status = ?", groupID, "active").
		Order("transacted_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTransactionsByUser(userID uuid.UUID) ([]models.Transaction, error) {
	db := database.DB

	var transactions []models.Transaction
	if err := db.Where("paid_by = ? AND status = ?", userID, "active").
		Order("transacted_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func UpdateTransaction(id uuid.UUID, req dto.UpdateTransactionRequest) (*models.Transaction, error) {
	db := database.DB

	var transaction models.Transaction
	if err := db.First(&transaction, "id = ? AND status = ?", id, "active").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	if req.Title != nil {
		transaction.Title = *req.Title
	}
	if req.Category != nil {
		transaction.Category = req.Category
	}
	if req.Notes != nil {
		transaction.Notes = req.Notes
	}
	if req.ReceiptURL != nil {
		transaction.ReceiptURL = req.ReceiptURL
	}

	if err := db.Save(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func DeleteTransaction(id uuid.UUID) error {
	db := database.DB

	result := db.Model(&models.Transaction{}).
		Where("id = ? AND status = ?", id, "active").
		Update("status", "deleted")

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

func GetGroupBalances(groupID uuid.UUID) ([]models.Balance, error) {
	db := database.DB

	var balances []models.Balance
	if err := db.Where("group_id = ?", groupID).Find(&balances).Error; err != nil {
		return nil, err
	}

	return balances, nil
}

// updateBalance is an internal helper that upserts the net balance between two users.
// It receives the active gorm.DB transaction so everything stays atomic.
func updateBalance(tx *gorm.DB, groupID, fromUser, toUser uuid.UUID, amount float64) error {
	var balance models.Balance

	err := tx.Where(
		"group_id = ? AND from_user_id = ? AND to_user_id = ?",
		groupID, fromUser, toUser,
	).First(&balance).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		balance = models.Balance{
			ID:         uuid.New(),
			GroupID:    groupID,
			FromUserID: fromUser,
			ToUserID:   toUser,
			NetAmount:  amount,
		}
		return tx.Create(&balance).Error
	}

	if err != nil {
		return err
	}

	balance.NetAmount += amount
	return tx.Save(&balance).Error
}
