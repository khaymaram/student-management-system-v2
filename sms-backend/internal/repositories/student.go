// Package repositories handles database access for student records.
// It encapsulates GORM operations and exposes CRUD/search/filter methods.
package repositories

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

type StudentRepository interface {
	GetAll() ([]models.Student, error)
	GetByID(id int) (*models.Student, error)
	Create(student *models.Student) error
	Update(student *models.Student) error
	Delete(id int) error

	Search(name string) ([]models.Student, error)
	FilterByGrade(grade int) ([]models.Student, error)
	GetHonors() ([]models.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{
		db: db,
	}
}

func (r *studentRepository) GetAll() ([]models.Student, error) {

	var students []models.Student

	err := r.db.Find(&students).Error

	return students, err
}

func (r *studentRepository) GetByID(id int) (*models.Student, error) {

	var student models.Student

	err := r.db.First(&student, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *studentRepository) Create(student *models.Student) error {
	// Insert a new student row into the database with GORM.
	// This is the lowest layer in the stack and is responsible for SQL persistence.
	return r.db.Create(student).Error
}

func (r *studentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) Delete(id int) error {
	// Remove the student row matching the provided ID.
	return r.db.Delete(&models.Student{}, "id = ?", id).Error
}

func (r *studentRepository) Search(name string) ([]models.Student, error) {

	var students []models.Student

	err := r.db.
		Where("name LIKE ?", "%"+name+"%").
		Find(&students).Error

	return students, err
}

func (r *studentRepository) FilterByGrade(grade int) ([]models.Student, error) {

	var students []models.Student

	err := r.db.
		Where("grade = ?", grade).
		Find(&students).Error

	return students, err
}

func (r *studentRepository) GetHonors() ([]models.Student, error) {

	var students []models.Student

	err := r.db.
		Where("gpa >= ?", 3.5).
		Find(&students).Error

	return students, err
}
