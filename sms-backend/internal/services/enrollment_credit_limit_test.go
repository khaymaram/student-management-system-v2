package services

import (
	"testing"

	"sms-backend/internal/dto"
	"sms-backend/internal/models"
	"sms-backend/internal/repositories"

	"gorm.io/gorm"
)

// ensure the service rejects enrollment that would exceed the 15-credit semester cap.
func TestEnrollmentService_EnrollRejectsOverCreditLimit(t *testing.T) {
	studentRepo := &stubStudentRepo{}
	courseRepo := &stubCourseRepo{course: models.Course{Code: "CSCI101", Credits: 6}}
	enrollmentRepo := &stubEnrollmentRepo{
		studentEnrollments: []models.Enrollment{
			{StudentID: 1, CourseCode: "MATH201", Course: &models.Course{Credits: 5}},
			{StudentID: 1, CourseCode: "BIO150", Course: &models.Course{Credits: 5}},
		},
	}

	service := NewEnrollmentService(enrollmentRepo, studentRepo, courseRepo)

	err := service.Enroll(1, dto.CreateEnrollmentRequest{CourseCode: "CSCI101"})
	if err == nil {
		t.Fatal("expected over-credit-limit enrollment to be rejected")
	}

	if err.Error() != "student cannot enroll in more than 15 credits in a semester" {
		t.Fatalf("unexpected error: %v", err)
	}
}

type stubStudentRepo struct{}

func (s *stubStudentRepo) GetAll() ([]models.Student, error) { return nil, nil }
func (s *stubStudentRepo) GetByID(id int) (*models.Student, error) {
	return &models.Student{ID: id, Name: "Test Student", Grade: 3, GPA: 3.5}, nil
}
func (s *stubStudentRepo) Create(student *models.Student) error { return nil }
func (s *stubStudentRepo) Update(student *models.Student) error { return nil }
func (s *stubStudentRepo) Delete(id int) error                  { return nil }
func (s *stubStudentRepo) Search(name string) ([]models.Student, error) {
	return nil, nil
}
func (s *stubStudentRepo) FilterByGrade(grade int) ([]models.Student, error) {
	return nil, nil
}
func (s *stubStudentRepo) GetHonors() ([]models.Student, error) { return nil, nil }

var _ repositories.StudentRepository = (*stubStudentRepo)(nil)

type stubCourseRepo struct {
	course models.Course
}

func (s *stubCourseRepo) GetByProfessor(professorId string) ([]models.Course, error) { return nil, nil }
func (s *stubCourseRepo) GetAll() ([]models.Course, error)                           { return nil, nil }
func (s *stubCourseRepo) GetByCode(code string) (*models.Course, error)              { return &s.course, nil }
func (s *stubCourseRepo) Create(course *models.Course) error                         { return nil }
func (s *stubCourseRepo) Update(course *models.Course) error                         { return nil }
func (s *stubCourseRepo) Delete(code string) error                                   { return nil }
func (s *stubCourseRepo) Search(title string) ([]models.Course, error)               { return nil, nil }
func (s *stubCourseRepo) FilterByCredits(credits int) ([]models.Course, error)       { return nil, nil }
func (s *stubCourseRepo) UnassignProfessor(professorId string) error                 { return nil }

var _ repositories.CourseRepository = (*stubCourseRepo)(nil)

type stubEnrollmentRepo struct {
	studentEnrollments []models.Enrollment
}

func (s *stubEnrollmentRepo) Create(enrollment *models.Enrollment) error { return nil }
func (s *stubEnrollmentRepo) Update(enrollment *models.Enrollment) error { return nil }
func (s *stubEnrollmentRepo) Delete(studentId int, courseCode string) error {
	return nil
}
func (s *stubEnrollmentRepo) Exists(studentId int, courseCode string) (bool, error) {
	return false, nil
}
func (s *stubEnrollmentRepo) GetByStudentAndCourse(studentId int, courseCode string) (*models.Enrollment, error) {
	return nil, gorm.ErrRecordNotFound
}
func (s *stubEnrollmentRepo) GetAll() ([]models.Enrollment, error) { return nil, nil }
func (s *stubEnrollmentRepo) GetByStudent(studentId int) ([]models.Enrollment, error) {
	return s.studentEnrollments, nil
}
func (s *stubEnrollmentRepo) GetByCourse(courseCode string) ([]models.Enrollment, error) {
	return nil, nil
}
func (s *stubEnrollmentRepo) DeleteByCourse(courseCode string) error { return nil }
func (s *stubEnrollmentRepo) DeleteByStudent(studentId int) error    { return nil }

var _ repositories.EnrollmentRepository = (*stubEnrollmentRepo)(nil)
