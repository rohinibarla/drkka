package config

import (
	"os"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	DB     DBConfig
	Static StaticConfig
	CORS   CORSConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

// DBConfig holds database-related configuration
type DBConfig struct {
	Path string
}

// StaticConfig holds static file serving configuration
type StaticConfig struct {
	Dir string
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           getEnv("PORT", "8080"),
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
		DB: DBConfig{
			Path: getEnv("DB_PATH", "./drkka.db"),
		},
		Static: StaticConfig{
			Dir: getEnv("STATIC_DIR", "../frontend/"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:8080"),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
