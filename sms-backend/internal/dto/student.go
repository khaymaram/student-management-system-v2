// Package dto defines request payload structures for student API requests.
// These DTOs are used by Gin handlers to bind and validate JSON input.
package dto

// CreateStudentRequest is the JSON payload expected from the frontend when creating a student.
// GPA is intentionally omitted because it is calculated from course grades.
type CreateStudentRequest struct {
	StudentID   int     `json:"studentId" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Grade       int64   `json:"grade" binding:"required"`
	Scholarship float64 `json:"scholarship"`
	IsInState   bool    `json:"isInState"`
}

// UpdateStudentRequest allows updating student fields that are entered manually.
// GPA is intentionally omitted because it is calculated from course grades.
type UpdateStudentRequest struct {
	Name  string `json:"name"`
	Grade int64  `json:"grade"`
}