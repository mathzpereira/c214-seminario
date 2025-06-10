package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mathzpereira/c214-seminario/contact-list-api/routes"

	_ "github.com/mathzpereira/c214-seminario/contact-list-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Contact List API
// @version 1.0
// @description API para gerenciamento de lista de contatos
// @host localhost:8080
// @BasePath /

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
