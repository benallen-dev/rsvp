package main

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port   string
	DBPath string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./test-data.db"
	}

	return &Config{
		Port:   port,
		DBPath: dbPath,
	}
}
