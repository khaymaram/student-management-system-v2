// Package initializers provides startup helpers for the backend.
// Migrate runs schema generation for the student model.
package initializers

import (
	"gorm.io/gorm"

	"sms-backend/internal/models"
)

// Migrate applies the application schema changes using GORM auto-migration.
// It drops the existing local tables first so older incompatible definitions
// do not prevent the enrollment table from being created on SQL Server.
func Migrate(db *gorm.DB) error {
	if err := dropTableIfExists(db, &models.Enrollment{}); err != nil {
		return err
	}

	if err := dropTableIfExists(db, &models.Course{}); err != nil {
		return err
	}

	if err := dropTableIfExists(db, &models.Professor{}); err != nil {
		return err
	}

	if err := dropTableIfExists(db, &models.Finance{}); err != nil {
		return err
	}

	if err := dropTableIfExists(db, &models.Student{}); err != nil {
		return err
	}

	if err := dropTableIfExists(db, &models.Major{}); err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.Major{},
		&models.Student{},
		&models.Finance{},
		&models.Professor{},
		&models.Course{},
		&models.Enrollment{},
	)
}
func dropTableIfExists(db *gorm.DB, model interface{}) error {
	if db.Migrator().HasTable(model) {
		return db.Migrator().DropTable(model)
	}

	return nil
}
