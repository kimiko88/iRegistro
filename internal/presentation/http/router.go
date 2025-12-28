package http

import (
	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/middleware"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
	"github.com/k/iRegistro/internal/presentation/ws"
)

func NewRouter(authHandler *handlers.AuthHandler, wsHandler *ws.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	healthHandler := handlers.NewHealthHandler()

	// Public routes
	r.GET("/health", healthHandler.Check)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/password-reset", authHandler.RequestPasswordReset)
		auth.POST("/password-reset/confirm", authHandler.ResetPassword)

		// Protected routes
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware("your-secret-key")) // Should inject config
		{
			protected.POST("/2fa/enable", authHandler.Enable2FA)
			protected.POST("/2fa/verify", authHandler.Verify2FA)
		}
	}

	// WebSocket
	if wsHandler != nil {
		r.GET("/ws", wsHandler.ServeWS)
	}

	return r
}
