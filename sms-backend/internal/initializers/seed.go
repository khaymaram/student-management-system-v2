// Seed inserts sample student records into the database for local development.
package initializers

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

// Seed adds example student data for local development and testing.
func Seed(db *gorm.DB) error {
	students := []models.Student{
		{
			ID:    1122,
			Name:  "Happy Birthday",
			Grade: 1,
			
		},
	}
	courses := []models.Course{
		{
			Title: "Linear Algebra",
			Code: "MATH240",
			Credits: 4,
		},
		{
			Title: "Intro to OOP",
			Code: "CMSC131",
			Credits: 4,
		},
		{
			Title: "Calculus II",
			Code: "MATH141",
			Credits: 4,
		},
		{
			Title: "Organization of Programming Languages",
			Code: "CMSC330",
			Credits: 3,
		},
	}
	db.Create(&courses)
	return db.Create(&students).Error
}
