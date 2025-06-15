package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/hydr0g3nz/wallet_topup_system/config"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/controller"
	repository "github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/ent"
	usecase "github.com/hydr0g3nz/wallet_topup_system/internal/application"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/infra"
	"github.com/hydr0g3nz/wallet_topup_system/internal/infrastructure"
)

func main() {

	// cfg := config.LoadFromEnv()
	// dbCfg := NewDBConfigFromConfig(cfg)
	// _, err := infrastructure.ConnectDB(dbCfg)
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to connect to database: %v", err))
	// }
	// return
	fx.New(
		// Provide dependencies
		fx.Provide(
			config.LoadFromEnv,
			NewLoggerFromConfig,
			NewDBConfigFromConfig,
			infrastructure.ConnectDB,
			NewRedisConfigFromConfig,
			infrastructure.NewRedisClient,
			repository.NewUserRepository,
			// repository.NewTransactionRepository,
			// repository.NewWalletRepository,
			// repository.NewDBTransactionRepository,
			// NewWalletUsecaseWithInterfaces,
			usecase.NewUserUsecase,
			// controller.NewWalletController,
			controller.NewUserController,
			NewFiberServer,
		),

		// Invoke lifecycle hooks
		fx.Invoke(
			infrastructure.RunEntMigration,
			registerRoutes,
			startServer,
		),
	).Run()
}

// DatabaseParams wraps database-related dependencies
type DatabaseParams struct {
	fx.In
	DB     *gorm.DB // แก้ไขให้ตรงกับ type ที่ใช้จริง
	Logger *infrastructure.Logger
}

// ServerParams wraps server-related dependencies
type ServerParams struct {
	fx.In
	App    *fiber.App
	Config *config.Config
	Logger infra.Logger
}

// RouteParams wraps routing dependencies
type RouteParams struct {
	fx.In
	App            *fiber.App
	UserController *controller.UserController
}

// NewLoggerFromConfig creates logger from config
func NewLoggerFromConfig(config *config.Config) (infra.Logger, error) {
	return infrastructure.NewLogger(config.IsProduction())
}

// NewDBConfigFromConfig extracts database config from main config
func NewDBConfigFromConfig(config *config.Config) *infrastructure.DBConfig {
	return &config.Database
}

// NewRedisConfigFromConfig extracts redis config from main config
func NewRedisConfigFromConfig(config *config.Config) infrastructure.CacheConfig {
	return config.Cache
}

// // WalletUsecaseParams wraps all dependencies needed for WalletUsecase
// type WalletUsecaseParams struct {
// 	fx.In
// 	UserRepo        *repository.UserRepository
// 	TransactionRepo *repository.TransactionRepository
// 	WalletRepo      *repository.WalletRepository
// 	Cache           *infrastructure.RedisClient
// 	DBTxRepo        *repository.DBTransactionRepository
// 	Logger          *infrastructure.Logger
// 	Config          *config.Config
// }

// // NewWalletUsecaseWithInterfaces creates WalletUsecase with proper interface conversion
// func NewWalletUsecaseWithInterfaces(params WalletUsecaseParams) usecase.WalletUsecase {
// 	// หากจำเป็นต้องแปลงเป็น interface ให้ทำแบบนี้:
// 	// var userRepo repository.UserRepository = params.UserRepo
// 	// var transactionRepo repository.TransactionRepository = params.TransactionRepo
// 	// var walletRepo repository.WalletRepository = params.WalletRepo
// 	// var cache infra.CacheService = params.Cache
// 	// var dbTxRepo repository.DBTransaction = params.DBTxRepo
// 	// var logger infra.Logger = params.Logger

// 	return usecase.NewWalletUsecase(
// 		params.UserRepo,
// 		params.TransactionRepo,
// 		params.WalletRepo,
// 		params.Cache,
// 		params.DBTxRepo,
// 		params.Logger,
// 		*params.Config,
// 	)
// }

// NewFiberServer creates a new Fiber server instance
func NewFiberServer(config *config.Config) *fiber.App {
	return infrastructure.NewFiber(infrastructure.ServerConfig{
		Address:      config.Server.Port,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
	})
}

// // runMigrations runs database migrations during startup
// func runMigrations(lc fx.Lifecycle, params DatabaseParams) {
// 	lc.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {
// 			params.Logger.Info("Running database migrations", nil)
// 			if err := infrastructure.MigrateDB(params.DB); err != nil {
// 				params.Logger.Error("Failed to run database migrations", map[string]interface{}{
// 					"error": err.Error(),
// 				})
// 				return err
// 			}
// 			params.Logger.Info("Database migrations completed successfully", nil)
// 			return nil
// 		},
// 	})
// }

// // seedDatabase seeds the database with initial data during startup
// func seedDatabase(lc fx.Lifecycle, params DatabaseParams) {
// 	lc.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {
// 			params.Logger.Info("Seeding database with initial data", nil)
// 			if err := infrastructure.SeedDB(params.DB); err != nil {
// 				params.Logger.Error("Failed to seed database", map[string]interface{}{
// 					"error": err.Error(),
// 				})
// 				return err
// 			}
// 			params.Logger.Info("Database seeding completed successfully", nil)
// 			return nil
// 		},
// 	})
// }

// registerRoutes registers all API routes
func registerRoutes(params RouteParams) {
	// Setup API routes
	api := params.App.Group("/api/v1")
	params.UserController.RegisterRoutes(api)
}

// startServer starts the Fiber server
func startServer(lc fx.Lifecycle, params ServerParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				params.Logger.Info("Starting server", map[string]interface{}{
					"port": params.Config.Server.Port,
				})

				if err := params.App.Listen(fmt.Sprintf(":%s", params.Config.Server.Port)); err != nil {
					params.Logger.Fatal("Failed to start server", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Shutting down server", nil)
			return params.App.Shutdown()
		},
	})
}
