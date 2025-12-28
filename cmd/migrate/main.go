package main

import (
	"log"

	"github.com/k/iRegistro/internal/config"
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
	// AutoMigrate will be added here as domains are implemented
	// if err := db.AutoMigrate(&domain.User{}, ...); err != nil { ... }

	// Create a dummy table check or similar if needed for now,
	// but mostly this is a placeholder for future migrations.
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get generic db interface", zap.Error(err))
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("Migrations completed successfully (connected and ready).")
}
