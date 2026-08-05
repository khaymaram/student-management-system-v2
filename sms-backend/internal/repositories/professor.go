package repositories

import (
	"sms-backend/internal/models"

	"gorm.io/gorm"
)

type ProfessorRepository interface {
	GetAll() ([]models.Professor, error)
	GetByID(id string) (*models.Professor, error)
	GetPaginated(page int, limit int, professorId string, name string)([]models.Professor, int64, error)
	Create(professor *models.Professor) error
	Update(professor *models.Professor) error
	Delete(id string) error
	Search(name string) ([]models.Professor, error)
}

type professorRepository struct {
	db *gorm.DB
}

func NewProfessorRepository(db *gorm.DB) ProfessorRepository {
	return &professorRepository{
		db: db,
	}
}

func (r *professorRepository) GetAll() ([]models.Professor, error) {

	var professors []models.Professor

	err := r.db.Find(&professors).Error
	r.db.Preload("Courses")

	return professors, err
}

func (r *professorRepository) GetPaginated(page int, limit int, professorId string, name string) ([]models.Professor, int64, error){
	var professors []models.Professor
	var total int64
	query := r.db.Model(&models.Professor{})

	if professorId != ""{
		query = query.Where(
			"LOWER(id) = LOWER(?)",
			professorId,
		)
	}

	if name != "" {
		query = query.Where(
			"name LIKE ?",
			"%"+name+"%",
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.
		Order("id").
		Offset(offset).
		Limit(limit).
		Find(&professors).Error; err != nil {
			return nil, 0, err
		}
	return professors, total, nil
}

func (r *professorRepository) GetByID(id string) (*models.Professor, error) {

	var professor models.Professor

	err := r.db.First(&professor, "LOWER(id) = LOWER(?)", id).Error
	r.db.Preload("Courses")

	if err != nil {
		return nil, err
	}

	return &professor, nil
}

func (r *professorRepository) Create(professor *models.Professor) error {
	// Insert a new student row into the database with GORM.
	// This is the lowest layer in the stack and is responsible for SQL persistence.
	return r.db.Create(professor).Error
}

func (r *professorRepository) Update(professor *models.Professor) error {
	return r.db.Save(professor).Error
}

func (r *professorRepository) Delete(id string) error {
	return r.db.Delete(&models.Professor{}, "LOWER(id) = LOWER(?)", id).Error
}

func (r *professorRepository) Search(name string) ([]models.Professor, error) {

	var professors []models.Professor

	err := r.db.
		Where("name LIKE ?", "%"+name+"%").
		Find(&professors).Error

	return professors, err
}
