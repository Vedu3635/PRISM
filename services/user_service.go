package services

import (
	"time"

	"github.com/Vedu3635/PRISM.git/database"
	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
