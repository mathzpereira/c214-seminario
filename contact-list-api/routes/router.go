package routes

import (
	"github.com/mathzpereira/c214-seminario/contact-list-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	contactGroup := router.Group("/contacts")
	{
		contactGroup.GET("/", handlers.GetContacts)
		contactGroup.POST("/", handlers.CreateContact)
		contactGroup.GET("/:id", handlers.GetContactByID)
		contactGroup.PUT("/:id", handlers.UpdateContactById)
		contactGroup.DELETE("/:id", handlers.DeleteContact)
		contactGroup.GET("/summary", handlers.GetContactsSummary)
		contactGroup.GET("/search", handlers.SearchContactsByName)
		contactGroup.GET("/email-providers", handlers.GetEmailProviders)
	}
}
