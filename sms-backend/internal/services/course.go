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
}

func NewCourseService(repo repositories.CourseRepository) CourseService {

	return &courseService{
		repository: repo,
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
		// ID:    uint(req.CourseID),
		Title:  req.Title,
		Code: req.Code,
		Credits:   req.Credits,
	}

	return c.repository.Create(&course)
}

func (c *courseService) Update(code string, req dto.UpdateCourseRequest) error {
	course, err := c.repository.GetByCode(code)
	if err != nil {
		return err
	}

	// var course *models.Course
	// for i := range courses {
	// 	if int(courses[i].ID) == id {
	// 		course = &courses[i]
	// 		break
	// 	}
	// }

	// if course == nil {
	// 	return gorm.ErrRecordNotFound
	// }

	if req.Title != "" {
		course.Title = req.Title
	}

	// if req.Code != "" {
	// 	course.Code = req.Code
	// }

	if req.Credits != 0 {
		course.Credits = req.Credits
	}

	return c.repository.Update(course)
}

func (c *courseService) Delete(code string) error {
	return c.repository.Delete(code)
}

func (c *courseService) Search(title string) ([]models.Course, error) {
	return c.repository.Search(title)
}

func (c *courseService) FilterByCredits(credits int) ([]models.Course, error) {
	return c.repository.FilterByCredits(credits)
}

