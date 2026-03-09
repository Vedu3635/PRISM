package routes

import (
	"github.com/Vedu3635/PRISM.git/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	api := router.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Users
	users := api.Group("/users")
	{
		users.POST("/", handlers.CreateUser)
		users.GET("/", handlers.GetUsers)
		users.GET("/:id", handlers.GetUserByID)
	}

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

	api.POST("/transactions", handlers.CreateTransaction)
	// router.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "PRISM backend running")
	// })

	// router.GET("/health", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"status": "ok",
	// 	})
	// })
}
