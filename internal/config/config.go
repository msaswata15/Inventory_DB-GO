package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration derived from environment variables.
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	ServerPort string
	SMTPHost   string
	SMTPPort   int
	SMTPUser   string
	SMTPPass   string
	SMTPTo     string
}

// Load reads configuration from the environment, optionally loading a .env file first.
func Load() (*Config, error) {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "inventory_db"),
		SSLMode:    getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		SMTPHost:   getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:   mustAtoi(getEnv("SMTP_PORT", "587")),
		SMTPUser:   getEnv("SMTP_USER", "your@gmail.com"),
		SMTPPass:   getEnv("SMTP_PASSWORD", "app-password"),
		SMTPTo:     getEnv("SMTP_TO", "receiver@gmail.com"),
	}
	return cfg, nil
}

func mustAtoi(v string) int {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
