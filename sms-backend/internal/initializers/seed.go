// Seed inserts sample student records into the database for local development.
package initializers

import (
	"sms-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed adds example student data for local development and testing.
func Seed(db *gorm.DB) error {
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
	studentHash, _ := bcrypt.GenerateFromPassword([]byte("Student1122!"), bcrypt.DefaultCost)
	teacherHash, _ := bcrypt.GenerateFromPassword([]byte("TeacherP1001!"), bcrypt.DefaultCost)
	adminSubject, studentSubject, teacherSubject := "ADMIN001", "1122", "P1001"
	users := []models.User{
		{Name: "System Administrator", Email: "admin@grgi.edu", PasswordHash: string(adminHash), Role: "admin", SubjectID: &adminSubject},
		{Name: "Happy Birthday", Email: "happybirthday1122@grgi.edu", PasswordHash: string(studentHash), Role: "student", SubjectID: &studentSubject},
		{Name: "Mr Bean", Email: "mrbeanp1001@grgi.edu", PasswordHash: string(teacherHash), Role: "professor", SubjectID: &teacherSubject},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}
	majors := []models.Major{
		{ID: 1, Name: "Undeclared"},
		{ID: 2, Name: "Computer Science"},
		{ID: 3, Name: "Mathematics"},
		{ID: 4, Name: "Business"},
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
			Title:       "Intro to Psychology",
			Code:        "PSYC101",
			Credits:     3,
			MeetingDays: []string{"T", "Th"}, StartTime: "13:00",
		},
	}
	professors := []models.Professor{{ID: teacherSubject, Name: "Mr Bean"}}
	if err := db.Create(&professors).Error; err != nil {
		return err
	}
	courses[0].ProfessorID = &teacherSubject
	courses[1].ProfessorID = &teacherSubject
	if err := db.Create(&courses).Error; err != nil {
		return err
	}
	if err := db.Create(&students).Error; err != nil {
		return err
	}
	return db.Create(&models.Finance{
		StudentID:   1122,
		Tuition:     10000,
		Scholarship: 0,
		Paid:        0,
		IsInState:   true,
	}).Error
}
