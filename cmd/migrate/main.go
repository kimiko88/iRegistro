package main

import (
	"log"

	"github.com/k/iRegistro/internal/config"
	"github.com/k/iRegistro/internal/domain"
	"github.com/k/iRegistro/internal/infrastructure/persistence"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Connecting to database for migration...")
	db, err := persistence.NewDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	logger.Info("Running migrations...")
	if err := db.AutoMigrate(
		&domain.User{}, &domain.Session{}, &domain.RefreshToken{},
		// Communication
		&domain.Notification{}, &domain.NotificationPreference{},
		&domain.Conversation{}, &domain.Message{},
		&domain.ColloquiumSlot{}, &domain.ColloquiumBooking{},
		// Admin
		&domain.AuditLog{}, &domain.SchoolSettings{},
		&domain.UserImport{}, &domain.DataExport{},
	); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get generic db interface", zap.Error(err))
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("Migrations completed successfully (connected and ready).")
}
