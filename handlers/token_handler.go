package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
)

// GetToken godoc
//
//	@Summary		Get Firebase ID token
//	@Description	Exchanges email and password for a Firebase ID token. Use the returned idToken as the Bearer token for all protected endpoints.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginRequest		true	"Email and password"
//	@Success		200		{object}	dto.TokenResponse		"Firebase ID token and metadata"
//	@Failure		400		{object}	map[string]string		"validation error"
//	@Failure		401		{object}	map[string]string		"invalid credentials"
//	@Failure		500		{object}	map[string]string		"internal server error"
//	@Router			/auth/token [post]
func GetToken(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.GetFirebaseToken(req.Email, req.Password)
	if err != nil {
		msg := err.Error()
		// Map known Firebase error messages to proper HTTP codes
		if msg == "EMAIL_NOT_FOUND" || msg == "INVALID_PASSWORD" || msg == "INVALID_LOGIN_CREDENTIALS" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		if msg == "USER_DISABLED" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "account has been disabled"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		IDToken:      result.IDToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		Email:        result.Email,
	})
}
