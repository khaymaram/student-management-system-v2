// Package services contains application business logic for enrollments.
// It validates that both the student and course exist, prevents duplicate
// enrollments, and orchestrates calls to the repository layer.
package services

import (
	"errors"

	"sms-backend/internal/dto"
	"sms-backend/internal/models"
	"sms-backend/internal/repositories"

	"gorm.io/gorm"
)

type EnrollmentService interface {
	Enroll(studentId int, req dto.CreateEnrollmentRequest) error
	Unenroll(studentId int, courseCode string) error
	UpdateGrade(studentId int, courseCode string, req dto.UpdateEnrollmentRequest) error

	GetByStudent(studentId int) ([]models.Enrollment, error)
	GetByCourse(courseCode string) ([]models.Enrollment, error)
}

type enrollmentService struct {
	repository        repositories.EnrollmentRepository
	studentRepository repositories.StudentRepository
	courseRepository  repositories.CourseRepository
}

func NewEnrollmentService(
	repo repositories.EnrollmentRepository,
	studentRepo repositories.StudentRepository,
	courseRepo repositories.CourseRepository,
) EnrollmentService {

	return &enrollmentService{
		repository:        repo,
		studentRepository: studentRepo,
		courseRepository:  courseRepo,
	}
}

func (s *enrollmentService) Enroll(studentId int, req dto.CreateEnrollmentRequest) error {
	// Confirm both sides of the relationship exist before creating the link,
	// so a bad studentId/courseCode fails with a clear message instead of a
	// raw foreign-key constraint error from the database.

	if _, err := s.studentRepository.GetByID(studentId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("student not found")
		}
		return err
	}

	if _, err := s.courseRepository.GetByCode(req.CourseCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	exists, err := s.repository.Exists(studentId, req.CourseCode)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("student is already enrolled in this course")
	}

	enrollment := models.Enrollment{
		StudentID:  studentId,
		CourseCode: req.CourseCode,
	}

	return s.repository.Create(&enrollment)
}

func (s *enrollmentService) Unenroll(studentId int, courseCode string) error {

	_, err := s.repository.GetByStudentAndCourse(studentId, courseCode)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return err
	}

	return s.repository.Delete(studentId, courseCode)
}

func (s *enrollmentService) UpdateGrade(studentId int, courseCode string, req dto.UpdateEnrollmentRequest) error {

	enrollment, err := s.repository.GetByStudentAndCourse(studentId, courseCode)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return err
	}

	enrollment.Grade = req.Grade

	return s.repository.Update(enrollment)
}

func (s *enrollmentService) GetByStudent(studentId int) ([]models.Enrollment, error) {
	return s.repository.GetByStudent(studentId)
}

func (s *enrollmentService) GetByCourse(courseCode string) ([]models.Enrollment, error) {
	return s.repository.GetByCourse(courseCode)
}
