package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port    int
	BaseURL string
}

type DatabaseConfig struct {
	ConnectionStr   string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type AuthConfig struct {
	TokenExpiration time.Duration
	JWTSecret       string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded")
	}

	//Server Config
	port, err := parseIntEnv("SERVER_PORT", "8080")
	if err != nil {
		return nil, err
	}
	baseURL := getStringEnv("SERVER_BASE_URL", "http://localhost:8080")

	//Database Config
	connectionStr := getStringEnv("DB_CONNECTION_STRING", "postgresql://postgres:postgres@localhost:5432/urlshortnerdb?sslmode=disable")
	maxOpenConns, err := parseIntEnv("DB_MAX_OPEN_CONNS", "25")
	if err != nil {
		return nil, err
	}
	maxIdleConns, err := parseIntEnv("DB_MAX_IDLE_CONNS", "25")
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := parseDurationEnv("DB_CONN_MAX_LIFETIME", "5m")
	if err != nil {
		return nil, err
	}

	//Auth Config
	tokenExpiration, err := parseDurationEnv("AUTH_TOKEN_EXPIRATION", "15m")
	if err != nil {
		return nil, err
	}
	jwtSecret := getStringEnv("AUTH_JWT_SECRET", "auth-jwt-secret-key")

	//Email Config
	smtpPort, _ := parseIntEnv("EMAIL_SMTP_PORT", "587")

	return &Config{
		Server: ServerConfig{
			Port:    port,
			BaseURL: baseURL,
		},
		Database: DatabaseConfig{
			ConnectionStr:   connectionStr,
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
		},
		Auth: AuthConfig{
			TokenExpiration: tokenExpiration,
			JWTSecret:       jwtSecret,
		},
		Email: EmailConfig{
			SMTPHost:     getStringEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
			SMTPPort:     smtpPort,
			SMTPUsername: getStringEnv("EMAIL_SMTP_USERNAME", ""),
			SMTPPassword: getStringEnv("EMAIL_SMTP_PASSWORD", ""),
			FromEmail:    getStringEnv("EMAIL_FROM", "noreply@securelink.com"),
		},
	}, nil
}

func getStringEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseIntEnv(key, defaultValue string) (int, error) {
	valueStr := getStringEnv(key, defaultValue)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", key, valueStr, err)
	}

	return value, nil
}

func parseDurationEnv(key, defaultValue string) (time.Duration, error) {
	valueStr := getStringEnv(key, defaultValue)
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", key, valueStr, err)
	}

	return value, nil
}
