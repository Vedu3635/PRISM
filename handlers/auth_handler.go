package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
)

// Signup godoc
//
//	@Summary		Register a new user
//	@Description	Creates a user in the DB and stamps a Firebase custom claim with the DB UUID.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.SignupRequest		true	"Signup payload"
//	@Success		201		{object}	map[string]interface{}	"created user"
//	@Failure		400		{object}	map[string]string		"validation error"
//	@Failure		409		{object}	map[string]string		"email or username already in use"
//	@Failure		500		{object}	map[string]string		"internal server error"
//	@Router			/auth/signup [post]
func Signup(c *gin.Context) {
	var req dto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.Signup(req)
	if err != nil {
		if err.Error() == "email already in use" || err.Error() == "username already taken" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// Login godoc
//
//	@Summary		Login with email and password
//	@Description	Validates credentials. Use the Firebase SDK separately to get an ID token.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginRequest		true	"Login payload"
//	@Success		200		{object}	map[string]interface{}	"user object"
//	@Failure		400		{object}	map[string]string		"validation error"
//	@Failure		401		{object}	map[string]string		"invalid credentials"
//	@Router			/auth/login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return the user — frontend uses Firebase SDK to get the ID token
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetMe godoc
//
//	@Summary		Get current user
//	@Description	Returns the profile of the authenticated user. UUID is extracted from Firebase token claims — no extra DB lookup.
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}	"current user"
//	@Failure		401	{object}	map[string]string		"unauthorized"
//	@Failure		404	{object}	map[string]string		"user not found"
//	@Failure		500	{object}	map[string]string		"internal server error"
//	@Router			/me [get]
func GetMe(c *gin.Context) {
	// user_id is set by AuthMiddleware — extracted from token claims, no DB hit
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
		return
	}

	user, err := services.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
