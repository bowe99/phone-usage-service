package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	LogLevel string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "phone_usage_db"),
			Timeout:  getDurationEnv("MONGO_TIMEOUT", 10*time.Second),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	if config.MongoDB.URI == "" {
		return nil, fmt.Errorf("MONGO_URI is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
