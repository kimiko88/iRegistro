package http

import (
	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/academic"
	"github.com/k/iRegistro/internal/application/admin"
	"github.com/k/iRegistro/internal/application/communication"
	"github.com/k/iRegistro/internal/application/director"
	"github.com/k/iRegistro/internal/application/reporting"
	"github.com/k/iRegistro/internal/application/secretary"
	"github.com/k/iRegistro/internal/domain"
	"github.com/k/iRegistro/internal/infrastructure/pdf"
	"github.com/k/iRegistro/internal/infrastructure/persistence"
	"github.com/k/iRegistro/internal/infrastructure/storage"
	"github.com/k/iRegistro/internal/middleware"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
	"github.com/k/iRegistro/internal/presentation/ws"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(authHandler *handlers.AuthHandler, wsHandler *ws.Handler, db *gorm.DB, hub *ws.Hub, logger *zap.Logger) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.PrometheusMiddleware())

	healthHandler := handlers.NewHealthHandler()

	// Public routes
	r.GET("/health", healthHandler.Check)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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

		// --- Service Initialization ---
		// 1. Communication (Core for others)
		commRepo := persistence.NewCommunicationRepository(db)
		notifService := communication.NewNotificationService(commRepo)

		// 2. Reporting (Uses Notification)
		reportingRepo := persistence.NewReportingRepository(db)
		pdfGen := pdf.NewMarotoGenerator()
		reportingService := reporting.NewReportingService(reportingRepo, pdfGen, notifService)
		reportingHandler := handlers.NewReportingHandler(reportingService)

		// --- Communication Module Setup ---
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

		// --- Admin Module Setup ---
		adminRepo := persistence.NewAdminRepository(db)
		auditService := admin.NewAuditService(adminRepo)
		adminService := admin.NewAdminService(adminRepo, userRepo, academicRepo, auditService) // Reuse academicRepo defined above
		importService := admin.NewUserImportService(adminRepo, userRepo, logger)
		exportService := admin.NewDataExportService(adminRepo)
		adminHandler := handlers.NewAdminHandler(adminService, auditService, importService, exportService)

		// Admin Routes
		// SuperAdmin
		sa := r.Group("/superadmin")
		sa.Use(middleware.AuthMiddleware("your-secret-key")) // Add Role check middleware
		{
			sa.POST("/schools", adminHandler.CreateSchool)
		}

		// School Admin
		adm := r.Group("/admin")
		adm.Use(middleware.AuthMiddleware("your-secret-key")) // Add Role check middleware
		{
			adm.GET("/settings", adminHandler.GetSettings)
			adm.PUT("/settings", adminHandler.UpdateSetting)
			adm.GET("/users", adminHandler.GetUsers)
			adm.GET("/audit-logs", adminHandler.GetAuditLogs)
			adm.POST("/data-export", adminHandler.RequestExport)
		}

		// Teacher Module
		teacherHandler := handlers.NewTeacherHandler(academicService)
		tch := r.Group("/teacher")
		tch.Use(middleware.AuthMiddleware("your-secret-key"))
		{
			tch.GET("/classes", teacherHandler.GetClasses)
			// Note: In real app, ensure teacher accesses only their classes
			tch.GET("/classes/:classId/students", teacherHandler.GetStudents)
			tch.GET("/classes/:classId/subjects/:subjectId/marks", teacherHandler.GetMarks)
			tch.POST("/marks", teacherHandler.CreateMark)
			tch.GET("/classes/:classId/absences", teacherHandler.GetAbsences)
			tch.POST("/classes/:classId/absences", teacherHandler.CreateAbsences)
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

		// --- Secretary Module Setup ---
		localStorage, _ := storage.NewLocalStorage("./uploads") // Simple local dir
		secService := secretary.NewSecretaryService(reportingRepo, pdfGen, localStorage, notifService)
		secHandler := handlers.NewSecretaryHandler(secService)

		directorService := director.NewDirectorService(academicRepo, reportingRepo)
		directorHandler := handlers.NewDirectorHandler(directorService)

		sec := r.Group("/secretary")
		sec.Use(middleware.AuthMiddleware("your-secret-key")) // Add Role check
		{
			sec.GET("/documents/inbox", secHandler.GetInbox)
			sec.GET("/documents/archive", secHandler.GetArchive)
			sec.POST("/documents/:id/approve", secHandler.ApproveDocument)
			sec.POST("/documents/:id/reject", secHandler.RejectDocument)
			sec.POST("/documents/:id/print-batch", secHandler.BatchPrint)
		}

		// --- Secretary Academic Management ---
		secAcademic := r.Group("/schools/:schoolId")
		secAcademic.Use(middleware.AuthMiddleware("your-secret-key"), middleware.RBACMiddleware(domain.RoleSecretary, domain.RoleAdmin))
		{
			// Class Management
			secAcademic.POST("/classes", academicHandler.CreateClass)
			secAcademic.POST("/assignments", academicHandler.AssignSubjectToClass)

			// Additional management if needed
			// secAcademic.POST("/students", academicHandler.CreateStudent)
			// secAcademic.POST("/enrollments", academicHandler.EnrollStudent)
		}

		// --- Files Setup ---
		fileHandler := handlers.NewFileHandler(localStorage)
		files := r.Group("/files")
		files.Use(middleware.AuthMiddleware("your-secret-key"))
		{
			files.GET("/download", fileHandler.DownloadFile)
		}

		// GraphQL
		r.POST("/query", handlers.GraphQLHandler(academicService, reportingService))
		r.GET("/playground", handlers.PlaygroundHandler())

		// Director routes
		directorRoutes := r.Group("/director")
		directorRoutes.Use(middleware.AuthMiddleware("your-secret-key"), middleware.RBACMiddleware(domain.RolePrincipal))
		{
			directorRoutes.GET("/kpi", directorHandler.GetKPIs)
			directorRoutes.GET("/documents/sign", directorHandler.GetDocumentsToSign)
			directorRoutes.POST("/documents/:id/sign", directorHandler.SignDocument)
		}

	}

	return r
}
