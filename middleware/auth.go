package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Vedu3635/PRISM.git/config"
)

// AuthMiddleware verifies the Firebase Bearer token and extracts db_user_id
// directly from custom claims — zero DB queries on every request.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid Authorization header",
			})
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Extract db_user_id from custom claims — set at signup, no DB query needed
		dbUserIDRaw, ok := token.Claims["db_user_id"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token missing db_user_id claim, please sign up first",
			})
			return
		}

		dbUserIDStr, ok := dbUserIDRaw.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid db_user_id claim format",
			})
			return
		}

		dbUserID, err := uuid.Parse(dbUserIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid db_user_id uuid in token",
			})
			return
		}

		// Make both available downstream
		c.Set("user_id", dbUserID) // uuid.UUID  — use this for DB queries
		c.Set("uid", token.UID)    // Firebase UID — use if needed
		// c.Set("email", token.Claims["email"])

		c.Next()
	}
}
