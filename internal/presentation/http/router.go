package http

import (
	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/academic"
	"github.com/k/iRegistro/internal/application/communication"
	"github.com/k/iRegistro/internal/application/reporting"
	"github.com/k/iRegistro/internal/infrastructure/pdf"
	"github.com/k/iRegistro/internal/infrastructure/persistence"
	"github.com/k/iRegistro/internal/middleware"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
	"github.com/k/iRegistro/internal/presentation/ws"
	"gorm.io/gorm"
)

func NewRouter(authHandler *handlers.AuthHandler, wsHandler *ws.Handler, db *gorm.DB, hub *ws.Hub) *gin.Engine {
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
		userRepo := persistence.NewUserRepository(db) // Reuse or create new
		academicRepo := persistence.NewAcademicRepository(db)
		broadcaster := ws.NewBroadcaster(hub) // hub is argument to NewRouter
		academicService := academic.NewAcademicService(academicRepo, userRepo, broadcaster)
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

		}

		// --- Reporting Module Setup ---
		reportingRepo := persistence.NewReportingRepository(db)
		pdfGen := pdf.NewMarotoGenerator()
		reportingService := reporting.NewReportingService(reportingRepo, pdfGen)
		reportingHandler := handlers.NewReportingHandler(reportingService)

		// --- Communication Module Setup ---
		commRepo := persistence.NewCommunicationRepository(db)
		notifService := communication.NewNotificationService(commRepo)
		msgService := communication.NewMessagingService(commRepo)
		colService := communication.NewColloquiumService(commRepo, notifService)
		commHandler := handlers.NewCommunicationHandler(notifService, msgService, colService)

		// Communication Routes
		comm := r.Group("/communication")
		comm.Use(middleware.AuthMiddleware("your-secret-key"))
		{
			// Notifications
			comm.GET("/notifications", commHandler.GetNotifications)
			comm.POST("/notifications/:id/read", commHandler.ReadNotification)

			// Messaging
			comm.POST("/conversations", commHandler.CreateConversation)
			comm.GET("/conversations", commHandler.GetConversations)
			comm.GET("/conversations/:id/messages", commHandler.GetMessages)
			comm.POST("/conversations/:id/messages", commHandler.SendMessage)

			// Colloquiums
			comm.POST("/slots", commHandler.CreateSlot) // Start simple, refine path usually /teachers/:id/slots
			comm.GET("/slots/available", commHandler.GetAvailableSlots)
			comm.POST("/bookings", commHandler.BookSlot)
		}

		// Route Group: /schools/:schoolId (Extensions)
		// Assuming we are within `schools` group context or similar, but structure above closes brackets.
		// Let's attach to `schools` group if possible, or create new.
		// The existing code closes `schools` group at line 67 in previous view.
		// Re-opening or adding inside the block is tricky with replace.
		// I will append new block for reporting after academic block.

		reporting := r.Group("/schools/:schoolId")
		reporting.Use(middleware.AuthMiddleware("your-secret-key"))
		{
			// Documents
			reporting.GET("/documents", reportingHandler.GetDocuments)
			reporting.POST("/classes/:classId/report-cards/generate", reportingHandler.GenerateReportCard)
			reporting.PATCH("/documents/:documentId/sign", reportingHandler.SignDocument)
			reporting.GET("/documents/:documentId/pdf", reportingHandler.GetDocumentPDF) // Fix Param access in handler if needed

			// PCTO
			reporting.GET("/pcto/projects", reportingHandler.GetPCTOProjects)
			reporting.POST("/pcto/projects", reportingHandler.CreatePCTOProject)

			// Orientation
			reporting.POST("/orientation/participations", reportingHandler.RegisterOrientation)
		}

		// GraphQL
		r.POST("/query", handlers.GraphQLHandler(academicService, reportingService))
		r.GET("/playground", handlers.PlaygroundHandler())
	}

	return r
}
