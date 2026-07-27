// Command seed populates the backend database with initial student records.
// It is useful for local development and testing.
package main

import (
	"log"

	"sms-backend/internal/initializers"
	"sms-backend/internal/pkg/config"
)

func main() {
	// Populate the database with starter records for development use.

	initializers.LoadEnv()

	cfg := config.Load()

	db, err := initializers.ConnectDB(cfg)

	if err != nil {
		log.Fatal(err)
	}

	err = initializers.Seed(db)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database seeded successfully.")
}
