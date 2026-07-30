// Package repositories handles database access for enrollment records.
// It encapsulates GORM operations and exposes CRUD/lookup methods for the
// many-to-many relationship between students and courses.
package repositories

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	Update(enrollment *models.Enrollment) error
	Delete(studentId int, courseCode string) error

	Exists(studentId int, courseCode string) (bool, error)
	GetByStudentAndCourse(studentId int, courseCode string) (*models.Enrollment, error)

	GetByStudent(studentId int) ([]models.Enrollment, error)
	GetByCourse(courseCode string) ([]models.Enrollment, error)

	DeleteByCourse(courseCode string) (error)
}

type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{
		db: db,
	}
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *enrollmentRepository) Update(enrollment *models.Enrollment) error {
	return r.db.Save(enrollment).Error
}

func (r *enrollmentRepository) Delete(studentId int, courseCode string) error {

	return r.db.
		Where("student_id = ? AND LOWER(course_code) = LOWER(?)", studentId, courseCode).
		Delete(&models.Enrollment{}).Error
}

func (r *enrollmentRepository) Exists(studentId int, courseCode string) (bool, error) {

	var count int64

	err := r.db.Model(&models.Enrollment{}).
		Where("student_id = ? AND LOWER(course_code) = LOWER(?)", studentId, courseCode).
		Count(&count).Error

	return count > 0, err
}

func (r *enrollmentRepository) GetByStudentAndCourse(studentId int, courseCode string) (*models.Enrollment, error) {

	var enrollment models.Enrollment

	err := r.db.
		Where("student_id = ? AND LOWER(course_code) = LOWER(?)", studentId, courseCode).
		First(&enrollment).Error

	if err != nil {
		return nil, err
	}

	return &enrollment, nil
}

// GetByStudent returns every course a student is enrolled in, with the
// related Course record preloaded so the frontend gets title/credits in
// one request.
func (r *enrollmentRepository) GetByStudent(studentId int) ([]models.Enrollment, error) {

	var enrollments []models.Enrollment

	err := r.db.
		Preload("Course").
		Where("student_id = ?", studentId).
		Find(&enrollments).Error

	return enrollments, err
}

// GetByCourse returns the roster of students enrolled in a course, with the
// related Student record preloaded.
func (r *enrollmentRepository) GetByCourse(courseCode string) ([]models.Enrollment, error) {

	var enrollments []models.Enrollment

	err := r.db.
		Preload("Student").
		Where("LOWER(course_code) = LOWER(?)", courseCode).
		Find(&enrollments).Error

	return enrollments, err
}

func (r *enrollmentRepository) DeleteByCourse(courseCode string) error {
	return r.db.Where("LOWER(course_code) = LOWER(?)", courseCode).
	Delete(&models.Enrollment{}).Error
}
