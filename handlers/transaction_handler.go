package handlers

import (
	"net/http"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTransaction godoc
//
//	@Summary		Create a transaction
//	@Description	Records a new expense in a group with participant splits.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateTransactionRequest	true	"Transaction payload"
//	@Success		201		{object}	models.Transaction				"created transaction"
//	@Failure		400		{object}	map[string]string				"validation error"
//	@Failure		500		{object}	map[string]string				"internal server error"
//	@Router			/transactions [post]
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

// GetTransactions godoc
//
//	@Summary		List all transactions
//	@Description	Returns all transactions accessible to the authenticated user.
//	@Tags			transactions
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.Transaction	"transactions"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/transactions [get]
func GetTransactions(c *gin.Context) {
	transactions, err := services.GetTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetTransactionByID godoc
//
//	@Summary		Get transaction by ID
//	@Description	Returns a single transaction by UUID.
//	@Tags			transactions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Transaction UUID"
//	@Success		200	{object}	models.Transaction	"transaction"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		404	{object}	map[string]string	"transaction not found"
//	@Router			/transactions/{id} [get]
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

// GetTransactionsByGroup godoc
//
//	@Summary		Get transactions for a group
//	@Description	Returns all transactions belonging to the specified group.
//	@Tags			groups
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Group UUID"
//	@Success		200	{array}		models.Transaction	"transactions"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/groups/{id}/transactions [get]
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

// GetTransactionsByUser godoc
//
//	@Summary		Get transactions for a user
//	@Description	Returns all transactions involving the specified user.
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"User UUID"
//	@Success		200	{array}		models.Transaction	"transactions"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/users/{id}/transactions [get]
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

// UpdateTransaction godoc
//
//	@Summary		Update a transaction
//	@Description	Updates mutable fields on a transaction — title, category, notes, receipt URL.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Transaction UUID"
//	@Param			body	body		dto.UpdateTransactionRequest	true	"Fields to update"
//	@Success		200		{object}	models.Transaction				"updated transaction"
//	@Failure		400		{object}	map[string]string				"invalid id or payload"
//	@Failure		500		{object}	map[string]string				"internal server error"
//	@Router			/transactions/{id} [put]
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

// DeleteTransaction godoc
//
//	@Summary		Delete a transaction
//	@Description	Soft-deletes a transaction by UUID.
//	@Tags			transactions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Transaction UUID"
//	@Success		200	{object}	map[string]string	"message"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/transactions/{id} [delete]
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
