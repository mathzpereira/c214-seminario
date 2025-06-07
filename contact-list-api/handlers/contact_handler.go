package handlers

import (
	"net/http"

	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
	"github.com/mathzpereira/c214-seminario/contact-list-api/services"

	"github.com/gin-gonic/gin"
)

func GetContacts(c *gin.Context) {
    contacts, err := services.GetAllContacts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, contacts)
}

func CreateContact(c *gin.Context) {
    var contact models.Contact
    if err := c.ShouldBindJSON(&contact); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := services.AddContact(contact); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, contact)
}
