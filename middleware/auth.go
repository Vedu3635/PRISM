package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Vedu3635/PRISM.git/config"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid Authorization header",
			})
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Set user info in context for handlers to use
		c.Set("uid", token.UID)
		c.Set("email", token.Claims["email"])

		c.Next()
	}
}
