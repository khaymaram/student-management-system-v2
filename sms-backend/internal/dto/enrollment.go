// Package dto defines request payload structures for enrollment API requests.
// These DTOs are used by Gin handlers to bind and validate JSON input.
package dto

// CreateEnrollmentRequest is the JSON payload sent when enrolling a student in a course.
type CreateEnrollmentRequest struct {
	CourseCode string `json:"courseCode" binding:"required"`
}

// UpdateEnrollmentRequest allows recording or changing the final grade for an enrollment.
type UpdateEnrollmentRequest struct {
	Grade string `json:"grade"`
}
