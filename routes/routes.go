package routes

import (
	"github.com/Vedu3635/PRISM.git/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	api := router.Group("/api")

	// Health
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Users
	users := api.Group("/users")
	{
		users.POST("/", handlers.CreateUser)
		users.GET("/", handlers.GetUsers)
		users.GET("/:id", handlers.GetUserByID)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}

	// Groups
	groups := api.Group("/groups")
	{
		groups.POST("/", handlers.CreateGroup)
		groups.GET("/", handlers.GetGroups)
		groups.GET("/:id", handlers.GetGroupsByID)
		groups.PUT("/:id", handlers.UpdateGroup)
		groups.DELETE("/:id", handlers.DeleteGroup)

		// Members
		members := groups.Group("/:id/members")
		{
			members.POST("/", handlers.AddMember)
			members.GET("/", handlers.GetGroupMembers)
			members.DELETE("/:memberID", handlers.RemoveMember)
		}

		groups.POST("/:id/leave", handlers.LeaveGroup)
	}

	// Transactions
	transactions := api.Group("/transactions")
	{
		transactions.POST("/", handlers.CreateTransaction)
	}
}
