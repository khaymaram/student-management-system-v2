// Package dto defines request payload structures for student API requests.
// These DTOs are used by Gin handlers to bind and validate JSON input.
package dto

// CreateStudentRequest is the JSON payload expected from the frontend when creating a student.
type CreateStudentRequest struct {
	StudentID int     `json:"studentId" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Grade     int64   `json:"grade" binding:"required"`
	GPA       float64 `json:"gpa" binding:"required"`
}

// UpdateStudentRequest allows partial updates for student fields.
type UpdateStudentRequest struct {
	Name  string  `json:"name"`
	Grade int64   `json:"grade"`
	GPA   float64 `json:"gpa"`
}
