package services

import (
	"errors"

	"sms-backend/internal/dto"
	"sms-backend/internal/models"
	"sms-backend/internal/repositories"

	"gorm.io/gorm"
)

type CourseService interface {
	GetAll() ([]models.Course, error)
	GetByCode(code string) (*models.Course, error)

	Create(dto.CreateCourseRequest) error

	Update(code string, req dto.UpdateCourseRequest) error

	Delete(code string) error

	Search(title string) ([]models.Course, error)

	FilterByCredits(credits int) ([]models.Course, error)
}

type courseService struct {
	repository repositories.CourseRepository
	enrollmentRepository repositories.EnrollmentRepository
}

func NewCourseService(repo repositories.CourseRepository, eRepo repositories.EnrollmentRepository) CourseService {

	return &courseService{
		repository: repo,
		enrollmentRepository: eRepo,
	}
}

func (c *courseService) GetAll() ([]models.Course, error) {
	return c.repository.GetAll()
}

func (c *courseService) GetByCode(code string) (*models.Course, error) {
	return c.repository.GetByCode(code)
}

func (c *courseService) Create(req dto.CreateCourseRequest) error {
	// Validate the request before saving it to the database.
	// This is the business-rule layer, where duplicate IDs and invalid GPA values are rejected.

	_, err := c.repository.GetByCode(req.Code)

	if err == nil {
		return errors.New("Course with that code already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if req.Credits < 0 {
		return errors.New("invalid number of credits")
	}

	course := models.Course{
		Title:       req.Title,
		Code:        req.Code,
		Credits:     req.Credits,
		ProfessorID: req.ProfessorID,
	}

	return c.repository.Create(&course)
}

func (c *courseService) Update(code string, req dto.UpdateCourseRequest) error {
	course, err := c.repository.GetByCode(code)
	if err != nil {
		return err
	}

	if req.Title != "" {
		course.Title = req.Title
	}

	if req.Credits != 0 {
		course.Credits = req.Credits
	}

	if req.ProfessorID != "" {
		course.ProfessorID = req.ProfessorID
	}
	return c.repository.Update(course)
}

func (c *courseService) Delete(code string) error {
	if _, err := c.repository.GetByCode(code); err != nil {
		return err
	}
	if err := c.enrollmentRepository.DeleteByCourse(code); err != nil {
		return err
	}
	return c.repository.Delete(code)
}

func (c *courseService) Search(title string) ([]models.Course, error) {
	return c.repository.Search(title)
}

func (c *courseService) FilterByCredits(credits int) ([]models.Course, error) {
	return c.repository.FilterByCredits(credits)
}


