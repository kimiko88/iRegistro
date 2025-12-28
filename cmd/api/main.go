package main

import (
	"log"

	"github.com/k/iRegistro/internal/config"
	"github.com/k/iRegistro/internal/infrastructure/logger"
	httpPresentation "github.com/k/iRegistro/internal/presentation/http"
	"go.uber.org/zap"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Setup Logger
	l, err := logger.New(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}
	defer l.Sync()
	zap.ReplaceGlobals(l)

	l.Info("Starting iRegistro", zap.String("env", cfg.Server.Mode))

	// 3. Setup Router
	r := httpPresentation.NewRouter()

	// 4. Start Server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		l.Fatal("Failed to start server", zap.Error(err))
	}
}
