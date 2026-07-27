// Package initializers provides startup setup helpers for the backend.
// It loads environment variables and database connections.
package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv reads the local .env file so configuration values are available at startup.
func LoadEnv() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(".env file not found")
	}
}
