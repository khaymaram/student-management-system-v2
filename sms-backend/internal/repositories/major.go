package repositories

import (
	"errors"
	"sms-backend/internal/models"
	"strings"

	"gorm.io/gorm"
)

type MajorRepository interface {
	GetAll() ([]models.Major, error)
	GetPaginated(page, limit int) ([]models.Major, int64, error)
	Exists(id int) (bool, error)
	Create(name string) (*models.Major, error)
	Delete(id int) error
}

func (r *majorRepository) Create(name string) (*models.Major, error) {
	major := &models.Major{Name: strings.TrimSpace(name)}
	if major.Name == "" {
		return nil, errors.New("major name is required")
	}
	if err := r.db.Create(major).Error; err != nil {
		return nil, err
	}
	return major, nil
}

func (r *majorRepository) Delete(id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var major models.Major
		if err := tx.First(&major, "id = ?", id).Error; err != nil {
			return err
		}
		if strings.EqualFold(major.Name, "Undeclared") {
			return errors.New("Undeclared cannot be deleted")
		}

		var undeclared models.Major
		if err := tx.Where("LOWER(name) = ?", "undeclared").First(&undeclared).Error; err != nil {
			return errors.New("Undeclared major is missing")
		}
		if err := tx.Model(&models.Student{}).Where("major_id = ?", id).Update("major_id", undeclared.ID).Error; err != nil {
			return err
		}
		return tx.Delete(&major).Error
	})
}

type majorRepository struct{ db *gorm.DB }

func NewMajorRepository(db *gorm.DB) MajorRepository {
	return &majorRepository{db: db}
}

func (r *majorRepository) GetAll() ([]models.Major, error) {
	var majors []models.Major
	err := r.db.Order("id").Find(&majors).Error
	return majors, err
}

func (r *majorRepository) GetPaginated(page, limit int) ([]models.Major, int64, error) {
	var majors []models.Major
	var total int64
	if err := r.db.Model(&models.Major{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Order("id").Offset((page - 1) * limit).Limit(limit).Find(&majors).Error
	return majors, total, err
}

func (r *majorRepository) Exists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Major{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
