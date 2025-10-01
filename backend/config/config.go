/*
Package config ...
*/
package config

import (
	"backend/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	API      APIConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Path         string
	WALMode      bool
	TimeoutSecs  int
	MaxOpenConns int
	MaxIdleConns int
}

type APIConfig struct {
	Key string
}

type RateLimit struct {
	PerIPLimit    int
	ClearInterval time.Duration
}

type ServerConfig struct {
	AllowedOrigins   []string
	ConnectionsLimit int
	RateLimit        RateLimit
	RequestSizeLimit int64
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	HandlerTimeout   time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Path:         utils.Must(getEnvWithoutDefault("DB_PATH")),
			WALMode:      getBoolEnv("DB_WAL_MODE", true),
			TimeoutSecs:  getIntEnv("DB_TIMEOUT", 30),
			MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 5),
			MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 2),
		},
		API: APIConfig{
			Key: utils.Must(getEnvWithoutDefault("API_KEY")),
		},
		Server: ServerConfig{
			AllowedOrigins:   getSliceEnv("ALLOWED_ORIGINS", []string{}),
			ConnectionsLimit: getIntEnv("CONNECTIONS_LIMIT", 100),
			RateLimit: RateLimit{
				PerIPLimit:    getIntEnv("RATE_LIMIT", 10),
				ClearInterval: time.Duration(getIntEnv("RATE_LIMIT_CLEAR_INTERVAL_SECS", 60)) * time.Second,
			},
			RequestSizeLimit: 1024 * 10, // 10 KB
			ReadTimeout:      10 * time.Second,
			WriteTimeout:     20 * time.Second,
			IdleTimeout:      60 * time.Second,
			HandlerTimeout:   25 * time.Second,
		},
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvWithoutDefault(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return value, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("Invalid boolean value for %s: %s", key, value)
		}
		return parsed
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Invalid integer value for %s: %s", key, value)
		}
		return parsed
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, len(parts))
		for i, part := range parts {
			result[i] = strings.TrimSpace(part)
		}

		return result
	}

	return defaultValue
}

func readFileFromEnvPath(envPath string) (string, error) {
	relativeFilePath := utils.Must(getEnvWithoutDefault(envPath))
	fileContents := utils.Must(os.ReadFile(relativeFilePath))
	return string(fileContents), nil
}
