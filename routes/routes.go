package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "PRISM backend running")
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
