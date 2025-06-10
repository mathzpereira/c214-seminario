package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
	"github.com/mathzpereira/c214-seminario/contact-list-api/services"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Error string `json:"error"`
}

// GetContacts godoc
// @Summary Lista todos os contatos
// @Tags Contacts
// @Produce json
// @Success 200 {array} models.Contact
// @Failure 500 {object} handlers.HTTPError
// @Router /contacts/ [get]
func GetContacts(c *gin.Context) {
	contacts, err := services.GetAllContacts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contacts)
}

// CreateContact godoc
// @Summary Cria um novo contato
// @Tags Contacts
// @Accept json
// @Produce json
// @Param contact body models.Contact true "Contato"
// @Success 201 {object} models.Contact
// @Failure 400 {object} handlers.HTTPError
// @Failure 500 {object} handlers.HTTPError
// @Router /contacts/ [post]
func CreateContact(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addedContact, err := services.AddContact(contact)
	if err != nil {
		if strings.Contains(err.Error(), "empty") || strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, addedContact)
}

// GetContactByID godoc
// @Summary Busca um contato por ID
// @Tags Contacts
// @Produce json
// @Param id path int true "ID do contato"
// @Success 200 {object} models.Contact
// @Failure 400,404 {object} handlers.HTTPError
// @Router /contacts/{id} [get]
func GetContactByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID deve ser número"})
		return
	}

	contact, err := services.GetContactByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contato não encontrado"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// UpdateContactById atualiza um contato existente por ID
// @Summary Atualiza um contato por ID
// @Description Atualiza os dados de um contato existente
// @Tags Contacts
// @Accept json
// @Produce json
// @Param id path int true "ID do contato"
// @Param contact body models.Contact true "Dados atualizados do contato"
// @Success 201 {object} models.Contact
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /contacts/{id} [put]
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
		if strings.Contains(err.Error(), "empty") || strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedContact)
}

// DeleteContact remove um contato por ID
// @Summary Remove um contato
// @Description Deleta um contato existente usando o ID
// @Tags Contacts
// @Param id path int true "ID do contato"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /contacts/{id} [delete]
func DeleteContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := services.DeleteContactById(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetContactsSummary retorna um resumo dos contatos
// @Summary Resumo dos contatos
// @Description Obtém estatísticas ou dados agregados sobre os contatos
// @Tags Contacts
// @Produce json
// @Success 200 {object} interface{} // pode substituir por um tipo exato se souber
// @Failure 500 {object} map[string]string
// @Router /contacts/summary [get]
func GetContactsSummary(c *gin.Context) {
	summary, err := services.GetContactsSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// SearchContactsByName busca contatos por nome
// @Summary Busca contatos
// @Description Busca contatos com base em parte do nome
// @Tags Contacts
// @Produce json
// @Param name query string true "Nome para busca parcial"
// @Success 200 {array} models.Contact
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /contacts/search [get]
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

// GetEmailProviders lista os provedores de e-mail dos contatos
// @Summary Lista provedores de e-mail
// @Description Retorna todos os domínios de e-mail utilizados pelos contatos
// @Tags Contacts
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} map[string]string
// @Router /contacts/email-providers [get]
func GetEmailProviders(c *gin.Context) {
	providers, err := services.GetEmailProviders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, providers)
}
