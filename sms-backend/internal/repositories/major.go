package repositories

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

type MajorRepository interface {
	GetAll() ([]models.Major, error)
	Exists(id int) (bool, error)
}

type majorRepository struct{ db *gorm.DB }

func NewMajorRepository(db *gorm.DB) MajorRepository {
	return &majorRepository{db: db}
}

func (r *majorRepository) GetAll() ([]models.Major, error) {
	var majors []models.Major
	err := r.db.Order("CASE WHEN name = 'Undeclared' THEN 0 ELSE 1 END, name").Find(&majors).Error
	return majors, err
}

func (r *majorRepository) Exists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Major{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
