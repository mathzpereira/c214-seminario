package handlers

import (
	"net/http"
	"strconv"

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

func GetContactByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": "ID deve ser número"})
		return
	}

	contact, err := services.GetContactByID(id)

	if err != nil {
		c.JSON(404, gin.H{"error": "Contato não encontrado"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func UpdateContactById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedContact, err := services.UpdateContactById(id, contact)

	if err != nil {
		c.JSON(404, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusCreated, updatedContact)
}

func DeleteContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := services.DeleteContactById(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.Status(http.StatusNoContent)
}

func GetContactsSummary(c *gin.Context) {
	summary, err := services.GetContactsSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func SearchContactsByName(c *gin.Context) {
	query := c.Query("name")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'name' is required"})
		return
	}

	contacts, err := services.SearchContactsByName(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contacts)
}