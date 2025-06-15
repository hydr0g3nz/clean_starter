package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/ent"
	_ "github.com/lib/pq"
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

// ConnectEntDB creates an ent.Client connected to PostgreSQL
func ConnectDB(cfg *DBConfig) (*ent.Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Try pinging the database
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return client, nil
}
func RunEntMigration(client *ent.Client) error {
	log.Println("Running ent migrations...")

	err := client.Schema.Create(context.Background())
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
		return err
	}

	log.Println("Ent migration completed")
	return nil
}
