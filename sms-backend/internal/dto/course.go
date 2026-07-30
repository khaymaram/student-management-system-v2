// Package dto defines request payload structures for student API requests.
// These DTOs are used by Gin handlers to bind and validate JSON input.
package dto

// CreateCourseRequest is the JSON payload expected from the frontend when creating a course.
type CreateCourseRequest struct {
	Code string `json:"code" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Credits     int     `json:"credits" binding:"required"`
	ProfessorID string `json:"professorId"`
}

// UpdateCourseRequest allows partial updates for course fields.
type UpdateCourseRequest struct {
	Title       string  `json:"title" `
	Code        string  `json:"code" `
	Credits     int     `json:"credits" `
	ProfessorID string `json:"professorId"`
}
