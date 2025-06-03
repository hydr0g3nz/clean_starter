package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConfig holds database connection configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectDB establishes a connection to the database
func ConnectDB(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
func MigrateDB(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto migrate all model
	err := db.AutoMigrate(
		&model.User{},
		&model.Wallet{},
		&model.Transaction{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedDB seeds the database with initial data
func SeedDB(db *gorm.DB) error {
	log.Println("Seeding database...")

	// Check if users already exist
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)

	// Only seed users if none exist
	if userCount == 0 {
		log.Println("Seeding users...")

		// Hash password - in production, use a stronger password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}

		// Create test users
		users := []model.User{
			{
				FirstName: "ธนาคาร",
				LastName:  "รุ่งเรือง",
				Email:     "tanakarn@example.com",
				Password:  string(hashedPassword),
				Phone:     "0812345678",
			},
			{
				FirstName: "สมหญิง",
				LastName:  "ใจดี",
				Email:     "somying@example.com",
				Password:  string(hashedPassword),
				Phone:     "0823456789",
			},
			{
				FirstName: "มานะ",
				LastName:  "มานี",
				Email:     "mana@example.com",
				Password:  string(hashedPassword),
				Phone:     "0834567890",
			},
			{
				FirstName: "ปิยะ",
				LastName:  "แสงทอง",
				Email:     "piya@example.com",
				Password:  string(hashedPassword),
				Phone:     "0845678901",
			},
			{
				FirstName: "วิชัย",
				LastName:  "เจริญ",
				Email:     "wichai@example.com",
				Password:  string(hashedPassword),
				Phone:     "0856789012",
			},
		}

		// Create users and their wallets in a transaction
		err = db.Transaction(func(tx *gorm.DB) error {
			for _, user := range users {
				// Create user
				if err := tx.Create(&user).Error; err != nil {
					log.Printf("Error seeding user: %v", err)
					return err
				}

				// Create wallet for user
				wallet := model.Wallet{
					Balance: 0.00, // Start with zero balance
				}

				if err := tx.Create(&wallet).Error; err != nil {
					log.Printf("Error creating wallet for user %s: %v", user.Email, err)
					return err
				}

				// Add example transaction for some users
				if user.ID%2 == 0 { // Only add transactions for even-numbered users
					// Add a completed transaction
					completedTx := model.Transaction{
						UserID:        user.ID,
						Amount:        500.00,
						PaymentMethod: "credit_card",
						Status:        "completed",
						ExpiresAt:     time.Now().Add(24 * time.Hour),
					}

					if err := tx.Create(&completedTx).Error; err != nil {
						log.Printf("Error creating transaction for user %s: %v", user.Email, err)
						return err
					}

					// Update wallet balance for completed transaction
					if err := tx.Model(&wallet).Update("balance", wallet.Balance+completedTx.Amount).Error; err != nil {
						log.Printf("Error updating wallet balance for user %s: %v", user.Email, err)
						return err
					}
				}

				// Add a pending transaction for the first user
				if user.ID == 1 {
					pendingTx := model.Transaction{
						UserID:        user.ID,
						Amount:        1000.00,
						PaymentMethod: "credit_card",
						Status:        "verified", // Pending verification
						ExpiresAt:     time.Now().Add(24 * time.Hour),
					}

					if err := tx.Create(&pendingTx).Error; err != nil {
						log.Printf("Error creating pending transaction for user %s: %v", user.Email, err)
						return err
					}
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("Error in transaction: %v", err)
			return err
		}

		log.Println("User and wallet seeding completed successfully")
	}

	log.Println("Database seeding completed successfully")
	return nil
}
