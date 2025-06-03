package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hydr0g3nz/wallet_topup_system/internal/infrastructure"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database infrastructure.DBConfig
	Cache    infrastructure.CacheConfig
	LogLevel string
	App      AppConfig
}
type AppConfig struct {
	MaxAcceptedAmount float64
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  int
	WriteTimeout int
	Host         string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret            string
	AccessExpiration  int // in minutes
	RefreshExpiration int // in hours
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("PORT", "8080"),
			// Environment:  getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),  // 10 seconds
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10), // 10 seconds
		},
		Database: infrastructure.DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "hospital_middleware"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Cache: infrastructure.CacheConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			Db:       getEnvAsInt("REDIS_DB", 0),
		},
		App: AppConfig{
			MaxAcceptedAmount: getEnvAsFloat("MAX_ACCEPTED_AMOUNT", 100000.0),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "release"
}

// getEnvAsInt gets an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return floatValue
		}
	}
	return defaultValue
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
