package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port    int
	BaseURL string
}

func Load() (*Config, error) {
	//Server Config
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			BaseURL: getEnv("SERVER_BASE_URL", "http://localhost:8080"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}