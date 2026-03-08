package routes

import (
	"github.com/Vedu3635/PRISM.git/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	api := router.Group("/api")

	api.POST("/users", handlers.CreateUser)
	api.GET("/users", handlers.GetUsers)
	api.GET("/users/:id", handlers.GetUserByID)

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
