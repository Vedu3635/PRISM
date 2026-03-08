package handlers

import (
	"net/http"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.CreateUser(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func GetUsers(c *gin.Context) {

	users, err := services.GetUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := services.GetUserByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
