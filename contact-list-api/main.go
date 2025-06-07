package main

import (
	"github.com/mathzpereira/c214-seminario/contact-list-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    routes.SetupRoutes(r)
    r.Run(":8080")
}