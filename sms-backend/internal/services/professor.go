package services

import (
	"errors"
	"gorm.io/gorm"
	"sms-backend/internal/dto"
	"sms-backend/internal/models"
	"sms-backend/internal/repositories"
)

type ProfessorService interface {
	GetAll() ([]models.Professor, error)

	GetPaginated(page int, limit int, professorId string, name string) ([]models.Professor, int64, error)

	GetByID(id string) (*models.Professor, error)

	Create(req dto.CreateProfessorRequest) error

	Update(id string, req dto.UpdateProfessorRequest) error

	Delete(id string) error

	Search(name string) ([]models.Professor, error)
}

type professorService struct {
	repository  repositories.ProfessorRepository
	courseRepo  repositories.CourseRepository
	authService *AuthService
}

func NewProfessorService(repo repositories.ProfessorRepository, courseRepo repositories.CourseRepository, auth *AuthService) ProfessorService {

	return &professorService{
		repository:  repo,
		courseRepo:  courseRepo,
		authService: auth,
	}
}

func (c *professorService) GetAll() ([]models.Professor, error) {
	return c.repository.GetAll()
}

func (c *professorService) GetPaginated(page int, limit int, professorId string, name string) ([]models.Professor, int64, error) {
	return c.repository.GetPaginated(page, limit, professorId, name)
}

func (c *professorService) GetByID(id string) (*models.Professor, error) {
	return c.repository.GetByID(id)
}

func (c *professorService) Create(req dto.CreateProfessorRequest) error {
	// Validate the request before saving it to the database.
	// This is the business-rule layer, where duplicate IDs and invalid GPA values are rejected.

	_, err := c.repository.GetByID(req.ID)

	if err == nil {
		return errors.New("Professor with that code already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	professor := models.Professor{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := c.repository.Create(&professor); err != nil {
		return err
	}
	if _, _, err := c.authService.CreateLinkedAccount(professor.Name, "professor", professor.ID); err != nil {
		_ = c.repository.Delete(professor.ID)
		return err
	}
	return nil
}

func (c *professorService) Update(id string, req dto.UpdateProfessorRequest) error {
	professor, err := c.repository.GetByID(id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		professor.Name = req.Name
	}

	return c.repository.Update(professor)
}

func (c *professorService) Delete(id string) error {
	if err := c.courseRepo.UnassignProfessor(id); err != nil {
		return err
	}
	if err := c.authService.DeleteLinkedAccount("professor", id); err != nil {
		return err
	}
	return c.repository.Delete(id)
}

func (c *professorService) Search(name string) ([]models.Professor, error) {
	return c.repository.Search(name)
}
