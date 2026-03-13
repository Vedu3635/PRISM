package handlers

import (
	"net/http"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := services.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func GetTransactions(c *gin.Context) {
	transactions, err := services.GetTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func GetTransactionByID(c *gin.Context) {
	id, err := parseTransactionID(c)
	if err != nil {
		return
	}

	tx, err := services.GetTransactionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func GetTransactionsByGroup(c *gin.Context) {
	groupIDParam := c.Param("id") // ← was "groupID", must match /:id/
	groupID, err := uuid.Parse(groupIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	transactions, err := services.GetTransactionsByGroup(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func GetTransactionsByUser(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	transactions, err := services.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func UpdateTransaction(c *gin.Context) {
	id, err := parseTransactionID(c)
	if err != nil {
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := services.UpdateTransaction(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func DeleteTransaction(c *gin.Context) {
	id, err := parseTransactionID(c)
	if err != nil {
		return
	}

	if err := services.DeleteTransaction(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted"})
}

// func GetGroupBalances(c *gin.Context) {
// 	groupIDParam := c.Param("groupID")
// 	groupID, err := uuid.Parse(groupIDParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
// 		return
// 	}

// 	balances, err := services.GetGroupBalances(groupID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, balances)
// }

func parseTransactionID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
	}
	return id, err
}
