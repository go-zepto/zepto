package utils

import (
	"log"
	"os"
)

// GetEnv get environment variable with fallback/default value
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// GetRequiredEnv get required environment variable
// If env var does not exist, throw error
func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalf("Required env var %s not defined\n", key)
	}
	return value
}
