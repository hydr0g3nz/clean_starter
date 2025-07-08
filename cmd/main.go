package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/config"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/controller"
	sqlcRepo "github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/sqlc"
	usecase "github.com/hydr0g3nz/wallet_topup_system/internal/application"
	"github.com/hydr0g3nz/wallet_topup_system/internal/infrastructure"
)

func main() {
	// Load configuration
	cfg := config.LoadFromEnv()

	// Setup logger
	logger, err := infrastructure.NewLogger(cfg.IsProduction())
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application")

	// Setup database
	db, err := infrastructure.ConnectDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer infrastructure.CloseDB(db)

	// Setup cache
	cache := infrastructure.NewRedisClient(cfg.Cache)
	defer cache.Close()

	// Setup repositories
	userRepo := sqlcRepo.NewUserRepository(db)

	// Setup use cases
	userUsecase := usecase.NewUserUsecase(userRepo, logger, cfg)

	// Setup controllers
	userController := controller.NewUserController(userUsecase)

	// Setup fiber server
	app := infrastructure.NewFiber(cfg)

	// Register routes
	api := app.Group("/api/v1")
	userController.RegisterRoutes(api)

	// Graceful shutdown
	go func() {
		logger.Info("Server starting", "port", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

// NewFiberServer creates a new Fiber server instance
func NewFiberServer(config *config.Config) *infrastructure.FiberApp {
	return infrastructure.NewFiber(infrastructure.ServerConfig{
		Address:      config.Server.Port,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
	})
}
