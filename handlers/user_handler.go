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

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Updates mutable fields — username, full name, phone, currency preference, password.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"User UUID"
//	@Param			body	body		dto.UpdateUserRequest	true	"Fields to update"
//	@Success		200		{object}	map[string]interface{}	"updated user"
//	@Failure		400		{object}	map[string]string		"invalid id or payload"
//	@Failure		500		{object}	map[string]string		"internal server error"
//	@Router			/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id, err := parseUserID(c)
	if err != nil {
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UpdateUser(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
//
//	@Summary		Soft-delete a user
//	@Description	Sets is_deleted = true. Does not remove the row.
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"User UUID"
//	@Success		200	{object}	map[string]string	"message"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, err := parseUserID(c)
	if err != nil {
		return
	}

	if err := services.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func parseUserID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
	}
	return id, err
}
