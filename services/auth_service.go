package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Vedu3635/PRISM.git/config"
	"github.com/Vedu3635/PRISM.git/database"
	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/models"
)

// Signup creates the user in DB, then stamps { db_user_id } as a
// Firebase custom claim so every future token carries the DB UUID — no
// extra DB lookup needed in the middleware.
func Signup(req dto.SignupRequest) (*models.User, error) {
	db := database.DB

	// Email uniqueness
	var existing models.User
	err := db.Where("email = ? AND is_deleted = ?", req.Email, false).First(&existing).Error
	if err == nil {
		return nil, errors.New("email already in use")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Username uniqueness
	err = db.Where("username = ? AND is_deleted = ?", req.Username, false).First(&existing).Error
	if err == nil {
		return nil, errors.New("username already taken")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:           uuid.New(),
		FirebaseUID:  req.FirebaseUID,
		Email:        req.Email,
		Username:     req.Username,
		FullName:     req.FullName,
		PasswordHash: string(hash),
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

	// Stamp the DB UUID as a custom claim on the Firebase user.
	// From this point every token the user gets will contain:
	//   "db_user_id": "<uuid>"
	claims := map[string]interface{}{
		"db_user_id": user.ID.String(),
	}
	if err := config.FirebaseAuth.SetCustomUserClaims(
		context.Background(), req.FirebaseUID, claims,
	); err != nil {
		// Roll back — don't leave a DB user without matching Firebase claims
		db.Delete(&user)
		return nil, errors.New("failed to set firebase claims: " + err.Error())
	}

	return &user, nil
}

// Login checks email + password and returns the user.
// The frontend should call Firebase SDK signInWithEmailAndPassword to get
// a fresh ID token (which will already contain db_user_id from signup).
func Login(req dto.LoginRequest) (*models.User, error) {
	db := database.DB

	var user models.User
	if err := db.Where("email = ? AND is_deleted = ?", req.Email, false).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// GetMe fetches the user by DB UUID (already extracted from token claims —
// no extra DB lookup in the hot path).
func GetMe(userID uuid.UUID) (*models.User, error) {
	db := database.DB

	var user models.User
	if err := db.First(&user, "id = ? AND is_deleted = ?", userID, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
