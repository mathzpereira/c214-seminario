package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mathzpereira/c214-seminario/contact-list-api/routes"
	"github.com/stretchr/testify/assert"
)

func TestGetContacts(t *testing.T) {
    router := gin.Default()
    routes.SetupRoutes(router)

    req, _ := http.NewRequest("GET", "/contacts/", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
