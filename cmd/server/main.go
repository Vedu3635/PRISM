package main

import (
	"log"

	"github.com/Vedu3635/PRISM.git/config"
	"github.com/Vedu3635/PRISM.git/database"
	"github.com/Vedu3635/PRISM.git/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	config.LoadEnv()
	config.InitFirebase()

	database.ConnectDB()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	routes.SetupRoutes(router)

	log.Println("Server running on port 8080")

	router.Run(":8080")
}
