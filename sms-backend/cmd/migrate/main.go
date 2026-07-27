// Command migrate runs database migrations for the backend.
// It uses GORM auto-migrate to create or update the student table.
package main

import (
	"log"

	"sms-backend/internal/initializers"
	"sms-backend/internal/pkg/config"
)

func main() {
	// Run database migrations so the student table exists before the app starts.

	initializers.LoadEnv()

	cfg := config.Load()

	db, err := initializers.ConnectDB(cfg)

	if err != nil {
		log.Fatal(err)
	}

	err = initializers.Migrate(db)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed successfully.")
}
