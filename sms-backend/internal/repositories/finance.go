package repositories

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

type financeRepository struct {
	db *gorm.DB
}

type FinanceRepository interface {
	GetAll() ([]models.Finance, error)
	GetPaginated(page int, limit int) ([]models.Finance, int64, error)
	Create(finance *models.Finance) error
	GetByStudentID(studentID int) (*models.Finance, error)
	Update(finance *models.Finance) error
	Delete(studentID int) error
}

func NewFinanceRepository(db *gorm.DB) FinanceRepository {
	return &financeRepository{
		db: db,
	}
}

func (r *financeRepository) GetPaginated(page int, limit int)([]models.Finance, int64, error){
	var finances []models.Finance
	var total int64
	query := r.db.Model(&models.Finance{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.Preload("Student").
		Order("student_id").
		Offset(offset).
		Limit(limit).
		Find(&finances).Error; err != nil {
			return nil, 0, err
		}
	return finances, total, nil
}

func (r *financeRepository) GetAll() ([]models.Finance, error) {
	var finances []models.Finance

	err := r.db.
		Preload("Student").
		Find(&finances).
		Error

	return finances, err
}

func (r *financeRepository) Create(finance *models.Finance) error {
	return r.db.Create(finance).Error
}

func (r *financeRepository) GetByStudentID(studentID int) (*models.Finance, error) {

	var finance models.Finance

	err := r.db.
		Preload("Student").
		First(&finance, "student_id = ?", studentID).
		Error

	return &finance, err
}

func (r *financeRepository) Update(finance *models.Finance) error {
	return r.db.Save(finance).Error
}

func (r *financeRepository) Delete(studentID int) error {
	return r.db.Delete(&models.Finance{}, "student_id = ?", studentID).Error
}
