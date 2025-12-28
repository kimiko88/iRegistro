package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	httpPresentation "github.com/k/iRegistro/internal/presentation/http"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Use the actual router implementation
	// For health check test, we don't need a real auth service
	authHandler := handlers.NewAuthHandler(nil)
	r = httpPresentation.NewRouter(authHandler, nil)

	// Perform Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())
}
