// Package initializers provides startup helpers for the backend.
// Migrate runs schema generation for the student model.
package initializers

import (
	"gorm.io/gorm"

	"sms-backend/internal/models"
)

// Migrate applies the student table schema changes using GORM auto-migration.
func Migrate(db *gorm.DB) error {

	return db.AutoMigrate(
		&models.Student{},
		&models.Course{},
	)
}
