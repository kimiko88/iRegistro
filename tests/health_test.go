package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	httpPresentation "github.com/k/iRegistro/internal/presentation/http"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	// Use the actual router implementation
	// For health check test, we don't need a real auth service
	authHandler := handlers.NewAuthHandler(nil)
	r := httpPresentation.NewRouter(authHandler, nil, nil, nil, zap.NewNop(), "test-secret")

	// Perform Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())
}
