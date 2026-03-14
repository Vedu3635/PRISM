package handlers

import (
	"net/http"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUser godoc
//
//	@Summary		Create a user
//	@Description	Creates a user record directly. For public registration use POST /auth/signup instead.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateUserRequest	true	"Create user payload"
//	@Success		201		{object}	map[string]interface{}	"created user"
//	@Failure		400		{object}	map[string]string		"validation error"
//	@Failure		500		{object}	map[string]string		"internal server error"
//	@Router			/users [post]
func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
//
//	@Summary		List all users
//	@Description	Returns all non-deleted users.
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}	"list of users"
//	@Failure		500	{object}	map[string]string		"internal server error"
//	@Router			/users [get]
func GetUsers(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
//
//	@Summary		Get user by ID
//	@Description	Returns a single user by UUID.
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string					true	"User UUID"
//	@Success		200	{object}	map[string]interface{}	"user"
//	@Failure		400	{object}	map[string]string		"invalid id"
//	@Failure		404	{object}	map[string]string		"user not found"
//	@Router			/users/{id} [get]
func GetUserByID(c *gin.Context) {
	id, err := parseUserID(c)
	if err != nil {
		return
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateMe godoc
//
//	@Summary		Update current user
//	@Description	Updates the authenticated user's own profile. ID is extracted from the Bearer token — no need to pass it in the URL.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.UpdateUserRequest	true	"Fields to update"
//	@Success		200		{object}	map[string]interface{}	"updated user"
//	@Failure		400		{object}	map[string]string		"validation error"
//	@Failure		401		{object}	map[string]string		"unauthorized"
//	@Failure		500		{object}	map[string]string		"internal server error"
//	@Router			/users/me [put]
func UpdateMe(c *gin.Context) {
	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UpdateUser(callerID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteMe godoc
//
//	@Summary		Delete current user
//	@Description	Soft-deletes the authenticated user's own account (is_deleted = true). ID is extracted from the Bearer token.
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]string	"message"
//	@Failure		401	{object}	map[string]string	"unauthorized"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/users/me [delete]
func DeleteMe(c *gin.Context) {
	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	if err := services.DeleteUser(callerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func parseUserID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
	}
	return id, err
}

// extractCallerID pulls the authenticated user's UUID from the gin context.
// It writes the appropriate error response and returns false if extraction fails.
// Reuse this anywhere you need the caller's ID from the token.
func extractCallerID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	id, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
		return uuid.Nil, false
	}

	return id, true
}
