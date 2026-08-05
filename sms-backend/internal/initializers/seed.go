// Seed inserts sample student records into the database for local development.
package initializers

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

// Seed adds example student data for local development and testing.
func Seed(db *gorm.DB) error {
	majors := []models.Major{
		{ID: 1, Name: "Undeclared"},
		{ID: 2, Name: "Computer Science"},
		{ID: 3, Name: "Mathematics"},
		{ID: 4, Name: "Business Administration"},
		{ID: 5, Name: "Biology"},
		{ID: 6, Name: "English"},
	}
	if err := db.Create(&majors).Error; err != nil {
		return err
	}
	students := []models.Student{
		{
			ID:      1122,
			Name:    "Happy Birthday",
			Grade:   1,
			MajorID: 1,
		},
	}
	courses := []models.Course{
		{
			Title:       "Linear Algebra",
			Code:        "MATH240",
			Credits:     4,
			MeetingDays: []string{"M", "W"}, StartTime: "08:00",
		},
		{
			Title:       "Intro to OOP",
			Code:        "CMSC131",
			Credits:     4,
			MeetingDays: []string{"T", "Th"}, StartTime: "09:30",
		},
		{
			Title:       "Calculus II",
			Code:        "MATH141",
			Credits:     4,
			MeetingDays: []string{"M", "W"}, StartTime: "11:00",
		},
		{
			Title:       "Organization of Programming Languages",
			Code:        "CMSC330",
			Credits:     3,
			MeetingDays: []string{"T", "Th"}, StartTime: "13:00",
		},
	}
	db.Create(&courses)
	return db.Create(&students).Error
}
