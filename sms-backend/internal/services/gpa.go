package services

import (
	"math"
	"strings"

	"sms-backend/internal/models"
	"sms-backend/internal/repositories"
)

func gradePoints(grade string) (float64, bool) {
	switch strings.ToUpper(strings.TrimSpace(grade)) {
	case "A":
		return 4.0, true
	case "B":
		return 3.0, true
	case "C":
		return 2.0, true
	case "D":
		return 1.0, true
	case "F":
		return 0.0, true
	default:
		return 0, false
	}
}

// calculateGPA returns nil when the student has no graded courses.
func calculateGPA(enrollments []models.Enrollment) *float64 {
	var qualityPoints float64
	var gradedCredits int

	for _, enrollment := range enrollments {
		if enrollment.Course == nil {
			continue
		}

		points, graded := gradePoints(enrollment.Grade)

		if !graded || enrollment.Course.Credits <= 0 {
			continue
		}

		qualityPoints += points * float64(enrollment.Course.Credits)
		gradedCredits += enrollment.Course.Credits
	}

	// No graded courses means there is no GPA.
	if gradedCredits == 0 {
		return nil
	}

	gpa := qualityPoints / float64(gradedCredits)

	// Round to two decimal places.
	gpa = math.Round(gpa*100) / 100

	return &gpa
}

func recalculateStudentGPA(
	studentID int,
	studentRepository repositories.StudentRepository,
	enrollmentRepository repositories.EnrollmentRepository,
) error {

	student, err := studentRepository.GetByID(studentID)
	if err != nil {
		return err
	}

	enrollments, err := enrollmentRepository.GetByStudent(studentID)
	if err != nil {
		return err
	}

	student.GPA = calculateGPA(enrollments)

	return studentRepository.Update(student)
}