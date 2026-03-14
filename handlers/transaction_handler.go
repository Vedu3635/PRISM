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
//	@Description	Records a new expense in a group with participant splits. The paid_by field must match the authenticated user.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateTransactionRequest	true	"Transaction payload"
//	@Success		201		{object}	models.Transaction				"created transaction"
//	@Failure		400		{object}	map[string]string				"validation error"
//	@Failure		403		{object}	map[string]string				"paid_by must match authenticated user"
//	@Failure		500		{object}	map[string]string				"internal server error"
//	@Router			/transactions [post]
func CreateTransaction(c *gin.Context) {
	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ── Force paid_by to the authenticated user — never trust client ──────────
	req.PaidBy = callerID

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
//	@Description	Returns all active transactions.
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
	groupID, err := uuid.Parse(c.Param("id"))
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
//	@Description	Returns all transactions where the specified user is the payer.
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"User UUID"
//	@Success		200	{array}		models.Transaction	"transactions"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/users/{id}/transactions [get]
func GetTransactionsByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
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
//	@Description	Updates mutable fields on a transaction. Only the user who paid can update it.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Transaction UUID"
//	@Param			body	body		dto.UpdateTransactionRequest	true	"Fields to update"
//	@Success		200		{object}	models.Transaction				"updated transaction"
//	@Failure		400		{object}	map[string]string				"invalid id or payload"
//	@Failure		403		{object}	map[string]string				"forbidden — only the payer can update"
//	@Failure		500		{object}	map[string]string				"internal server error"
//	@Router			/transactions/{id} [put]
func UpdateTransaction(c *gin.Context) {
	id, err := parseTransactionID(c)
	if err != nil {
		return
	}

	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	// ── Payer check ───────────────────────────────────────────────────────────
	if !requireTransactionPayer(c, id, callerID) {
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
//	@Description	Soft-deletes a transaction (status = deleted). Only the user who paid can delete it.
//	@Tags			transactions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Transaction UUID"
//	@Success		200	{object}	map[string]string	"message"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		403	{object}	map[string]string	"forbidden — only the payer can delete"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/transactions/{id} [delete]
func DeleteTransaction(c *gin.Context) {
	id, err := parseTransactionID(c)
	if err != nil {
		return
	}

	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	// ── Payer check ───────────────────────────────────────────────────────────
	if !requireTransactionPayer(c, id, callerID) {
		return
	}

	if err := services.DeleteTransaction(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted"})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func parseTransactionID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
	}
	return id, err
}

// requireTransactionPayer checks that the caller is the one who paid for the transaction.
// Writes 403 and returns false if not.
func requireTransactionPayer(c *gin.Context, txID uuid.UUID, callerID uuid.UUID) bool {
	isPayer, err := services.IsTransactionPayer(txID, callerID)
	if err != nil {
		if err.Error() == "transaction not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return false
	}
	if !isPayer {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the payer can modify this transaction"})
		return false
	}
	return true
}
