package services

import (
	"errors"
	"time"

	"github.com/Vedu3635/PRISM.git/database"
	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(req dto.CreateUserRequest) (*models.User, error) {
	db := database.DB

	var passwordHash string
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		passwordHash = string(hash)
	}

	user := models.User{
		ID:           uuid.New(),
		FirebaseUID:  req.FirebaseUID,
		Email:        req.Email,
		Username:     req.Username,
		FullName:     req.FullName,
		PasswordHash: passwordHash,
		Phone:        req.Phone,
		CurrencyPref: "INR",
		IsVerified:   false,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsers() ([]models.User, error) {
	db := database.DB

	var users []models.User
	if err := db.Where("is_deleted = ?", false).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(id uuid.UUID) (*models.User, error) {
	db := database.DB

	var user models.User
	if err := db.First(&user, "id = ? AND is_deleted = ?", id, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) (*models.User, error) {
	db := database.DB

	var user models.User
	if err := db.First(&user, "id = ? AND is_deleted = ?", id, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.CurrencyPref != nil {
		user.CurrencyPref = *req.CurrencyPref
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hash)
	}

	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(id uuid.UUID) error {
	db := database.DB

	result := db.Model(&models.User{}).Where("id = ? AND is_deleted = ?", id, false).Updates(map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
