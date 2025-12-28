package http

import (
	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/middleware"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	healthHandler := handlers.NewHealthHandler()

	// Basic health check
	r.GET("/health", healthHandler.Check)

	return r
}
