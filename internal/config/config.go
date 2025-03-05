package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
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
	Secret          string
}

func Load() (*Config, error) {
	//Server Config
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	//Database Config
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "25"))
	connMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))

	//Auth Config
	tokenExpiration, err := time.ParseDuration(getEnv("AUTH_TOKEN_EXPIRATION", "15m"))
	if err != nil {
		log.Fatal("Invalid AUTH_TOKEN_EXPIRATION format. Use a valid duration string like '15m'")
	}

	return &Config{
		Server: ServerConfig{
			Port:    port,
			BaseURL: getEnv("SERVER_BASE_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			ConnectionStr:   getEnv("DB_CONNECTION_STRING", "postgresql://postgres:postgres@localhost:5432/urlshortnerdb?sslmode=disable"),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
		},
		Auth: AuthConfig{
			TokenExpiration: tokenExpiration,
			Secret:          getEnv("AUTH_SECRET", "auth-secret-key"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
