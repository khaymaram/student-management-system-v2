// Package repositories handles database access for courses.
// It encapsulates GORM operations and exposes CRUD/search/filter methods.
package repositories

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

type CourseRepository interface {
	GetByProfessor(professorId string) ([]models.Course, error)
	GetAll() ([]models.Course, error)
	GetByCode(code string) (*models.Course, error)
	Create(student *models.Course) error
	Update(student *models.Course) error
	Delete(code string) error

	Search(title string) ([]models.Course, error)
	FilterByCredits(credits int) ([]models.Course, error)

	UnassignProfessor(professorId string) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{
		db: db,
	}
}
func (r *courseRepository) UnassignProfessor(professorId string) error {

	return r.db.
		Model(&models.Course{}).
		Where("professor_id = ?", professorId).
		Update("professor_id", nil).
		Error
}
func (r *courseRepository) GetByProfessor(professorId string) ([]models.Course, error) {

	var courses []models.Course

	err := r.db.Preload("Professor").
		Where("professor_id = ?", professorId).
		Find(&courses).Error

	return courses, err
}

func (r *courseRepository) GetAll() ([]models.Course, error) {

	var courses []models.Course

	err := r.db.Preload("Professor").Find(&courses).Error

	return courses, err
}

func (r *courseRepository) GetByCode(code string) (*models.Course, error) {

	var course models.Course

	err := r.db.Preload("Professor").First(&course, "LOWER(code) = LOWER(?)", code).Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) Create(course *models.Course) error {
	// Insert a new student row into the database with GORM.
	// This is the lowest layer in the stack and is responsible for SQL persistence.
	return r.db.Create(course).Error
}

func (r *courseRepository) Update(course *models.Course) error {
	updates := map[string]interface{}{
		"title":        course.Title,
		"credits":      course.Credits,
		"professor_id": course.ProfessorID,
	}

	return r.db.Model(&models.Course{}).
		Where("LOWER(code) = LOWER(?)", course.Code).
		Updates(updates).Error
}

func (r *courseRepository) Delete(code string) error {
	return r.db.Delete(&models.Course{}, "LOWER(code) = LOWER(?)", code).Error
}

func (r *courseRepository) Search(title string) ([]models.Course, error) {

	var courses []models.Course

	err := r.db.
		Where("title LIKE ?", "%"+title+"%").
		Find(&courses).Error

	return courses, err
}

func (r *courseRepository) FilterByCredits(credits int) ([]models.Course, error) {

	var courses []models.Course

	err := r.db.
		Where("credits = ?", credits).
		Find(&courses).Error

	return courses, err
}
