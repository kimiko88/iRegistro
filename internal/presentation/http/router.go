package http

import (
	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/academic"
	"github.com/k/iRegistro/internal/infrastructure/persistence"
	"github.com/k/iRegistro/internal/middleware"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
	"github.com/k/iRegistro/internal/presentation/ws"
	"gorm.io/gorm"
)

func NewRouter(authHandler *handlers.AuthHandler, wsHandler *ws.Handler, db *gorm.DB) *gin.Engine {
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

	// --- Academic Module Setup ---
	if db != nil {
		academicRepo := persistence.NewAcademicRepository(db)
		academicService := academic.NewAcademicService(academicRepo)
		academicHandler := handlers.NewAcademicHandler(academicService)

		// Route Group: /schools/:schoolId
		schools := r.Group("/schools/:schoolId")
		schools.Use(middleware.AuthMiddleware("your-secret-key")) // Reuse for now
		{
			schools.GET("/campuses", academicHandler.GetCampuses)
			schools.POST("/campuses", academicHandler.CreateCampus)

			schools.GET("/curriculums", academicHandler.GetCurriculums)
			schools.POST("/curriculums", academicHandler.CreateCurriculum)

			schools.GET("/classes", academicHandler.GetClasses)
			schools.POST("/classes", academicHandler.CreateClass)
			schools.GET("/classes/:classId", academicHandler.GetClassDetails)

			// Marks
			schools.POST("/classes/:classId/marks", academicHandler.CreateMark)
		}

		// GraphQL
		r.POST("/query", handlers.GraphQLHandler(academicService))
		r.GET("/playground", handlers.PlaygroundHandler())
	}

	return r
}
