package routes

import (
	"github.com/Vedu3635/PRISM.git/handlers"
	"github.com/Vedu3635/PRISM.git/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.Use(middleware.CORSMiddleware())

	api := router.Group("/api")

	// ─── Public ───────────────────────────────────────────────────────────────
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := api.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
	}

	// ─── Protected (Firebase token + db_user_id claim required) ───────────────

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{

		protected.GET("/me", handlers.GetMe)

		protected.GET("/debug/token", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetString("user_id"),
				"uid":     c.GetString("uid"),
				"email":   c.GetString("email"),
			})
		})
		// Users
		users := protected.Group("/users")
		{
			users.POST("/", handlers.CreateUser)
			users.GET("/", handlers.GetUsers)
			users.GET("/:id", handlers.GetUserByID)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
			users.GET("/:id/transactions", handlers.GetTransactionsByUser)
		}

		// Groups
		groups := protected.Group("/groups")
		{
			groups.POST("/", handlers.CreateGroup)
			groups.GET("/", handlers.GetGroups)
			groups.GET("/:id", handlers.GetGroupsByID)
			groups.PUT("/:id", handlers.UpdateGroup)
			groups.DELETE("/:id", handlers.DeleteGroup)
			groups.POST("/:id/leave", handlers.LeaveGroup)
			groups.GET("/:id/balances", handlers.GetGroupBalances)
			groups.GET("/:id/transactions", handlers.GetTransactionsByGroup)

			members := groups.Group("/:id/members")
			{
				members.POST("/", handlers.AddMember)
				members.GET("/", handlers.GetGroupMembers)
				members.DELETE("/:memberID", handlers.RemoveMember)
			}
		}

		// Transactions
		transactions := protected.Group("/transactions")
		{
			transactions.POST("/", handlers.CreateTransaction)
			transactions.GET("/", handlers.GetTransactions)
			transactions.GET("/:id", handlers.GetTransactionByID)
			transactions.PUT("/:id", handlers.UpdateTransaction)
			transactions.DELETE("/:id", handlers.DeleteTransaction)
		}
	}
}
