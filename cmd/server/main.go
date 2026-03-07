package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Vedu3635/PRISM/config"
	"github.com/Vedu3635/PRISM/routes"
)

func main() {

	config.LoadEnv()

	router := gin.Default()

	routes.SetupRoutes(router)

	log.Println("Server running on port 8080")

	router.Run(":8080")
}
