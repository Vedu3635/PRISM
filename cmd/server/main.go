// Package main is the entry point for the PRISM API.
//
//	@title			PRISM API
//	@version		1.0
//	@description	Expense splitting and group finance management backend.
//
//	@host		localhost:8080
//	@BasePath	/api
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Firebase ID token — format: "Bearer <token>"
package main

import (
	"log"

	"github.com/Vedu3635/PRISM.git/config"
	"github.com/Vedu3635/PRISM.git/database"
	_ "github.com/Vedu3635/PRISM.git/docs"
	"github.com/Vedu3635/PRISM.git/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	config.LoadEnv()
	config.InitFirebase()

	database.ConnectDB()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(router)

	log.Println("Server running on port 8080")
	log.Println("Swagger UI → http://localhost:8080/docs/index.html")

	router.Run(":8080")
}
