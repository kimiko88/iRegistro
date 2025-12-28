package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/k/iRegistro/internal/application/auth"
	"github.com/k/iRegistro/internal/config"
	"github.com/k/iRegistro/internal/infrastructure/logger"
	"github.com/k/iRegistro/internal/infrastructure/persistence"
	httpPresentation "github.com/k/iRegistro/internal/presentation/http"
	"github.com/k/iRegistro/internal/presentation/http/handlers"
	"github.com/k/iRegistro/internal/presentation/ws"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	_ = godotenv.Load() // Ignore error if .env doesn't exist

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

	// 3. Database
	db, err := persistence.NewDB(cfg.Database)
	if err != nil {
		l.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 4. Setup Dependencies
	userRepo := persistence.NewUserRepository(db)
	authRepo := persistence.NewAuthRepository(db)
	authService := auth.NewAuthService(
		userRepo,
		authRepo,
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessDuration,
		cfg.Auth.RefreshDuration,
	)
	authHandler := handlers.NewAuthHandler(authService)

	// WebSocket
	hub := ws.NewHub()
	go hub.Run()
	wsHandler := ws.NewHandler(hub, cfg.Auth.JWTSecret)

	// 5. Setup Router
	r := httpPresentation.NewRouter(authHandler, wsHandler, db, hub, l, cfg.Auth.JWTSecret)

	// 6. Start Server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		l.Fatal("Failed to start server", zap.Error(err))
	}
}
