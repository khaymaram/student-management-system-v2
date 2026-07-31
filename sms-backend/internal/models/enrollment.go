package models

import "time"

// Enrollment represents one student's registration in one course.
// A (student_id, course_code) pair is unique — a student cannot be
// enrolled in the same course more than once.
type Enrollment struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	StudentID int `gorm:"column:student_id;type:int;not null;uniqueIndex:idx_student_course" json:"studentId"`

	CourseCode string `gorm:"column:course_code;type:nvarchar(256);not null;uniqueIndex:idx_student_course" json:"courseCode"`

	// Grade is the letter grade earned in the course. It stays empty
	// until an instructor records a final grade for the student.
	Grade string `gorm:"column:grade" json:"grade"`

	CreatedAt time.Time `gorm:"column:created_at" json:"enrolledAt"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`

	// Student and Course are populated via Preload so API responses can
	// include the related record without a second round-trip from the frontend.
	Student *Student `gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE" json:"student,omitempty"`
	Course  *Course  `gorm:"foreignKey:CourseCode;references:Code;constraint:OnDelete:CASCADE" json:"course,omitempty"`
}

func (Enrollment) TableName() string {
	return "dbo.Enrollments"
}
