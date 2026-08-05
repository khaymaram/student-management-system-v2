package services

import (
	"errors"
	"fmt"

	"sms-backend/internal/dto"
	"sms-backend/internal/models"
	"sms-backend/internal/repositories"

	"gorm.io/gorm"
)

type CourseService interface {
	GetAll() ([]models.Course, error)
	GetByCode(code string) (*models.Course, error)
	GetByProfessor(professorId string) ([]models.Course, error)

	GetPaginated(page int, limit int, code string, title string, credits *int, professorId string) ([]models.Course, int64, error)

	Create(dto.CreateCourseRequest) error

	Update(code string, req dto.UpdateCourseRequest) error

	Delete(code string) error

	Search(title string) ([]models.Course, error)

	FilterByCredits(credits int) ([]models.Course, error)
}

type courseService struct {
	repository           repositories.CourseRepository
	enrollmentRepository repositories.EnrollmentRepository
	studentRepository    repositories.StudentRepository
}

func NewCourseService(repo repositories.CourseRepository, eRepo repositories.EnrollmentRepository, sRepo repositories.StudentRepository) CourseService {

	return &courseService{
		repository:           repo,
		enrollmentRepository: eRepo,
		studentRepository:    sRepo,
	}
}

func (c *courseService) GetPaginated(
	page int, limit int, code string, title string, credits *int, professorId string,
) ([]models.Course, int64, error) {
	return c.repository.GetPaginated(page, limit, code, title, credits, professorId)
}

func (c *courseService) GetByProfessor(professorId string) ([]models.Course, error) {
	return c.repository.GetByProfessor(professorId)
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

	days, startTime, err := normalizeSchedule(req.MeetingDays, req.StartTime)
	if err != nil {
		return err
	}

	course := models.Course{
		Title:       req.Title,
		Code:        req.Code,
		Credits:     req.Credits,
		ProfessorID: req.ProfessorID,
		MeetingDays: days,
		StartTime:   startTime,
	}
	if err := c.ensureProfessorAvailable(&course); err != nil {
		return err
	}

	return c.repository.Create(&course)
}

func (c *courseService) Update(code string, req dto.UpdateCourseRequest) error {
	course, err := c.repository.GetByCode(code)
	if err != nil {
		return err
	}
	oldCredits := course.Credits

	if req.Title != "" {
		course.Title = req.Title
	}

	if req.Credits != 0 {
		course.Credits = req.Credits
	}

	course.ProfessorID = req.ProfessorID

	days, startTime, err := normalizeSchedule(req.MeetingDays, req.StartTime)
	if err != nil {
		return err
	}
	course.MeetingDays = days
	course.StartTime = startTime

	if err := c.ensureProfessorAvailable(course); err != nil {
		return err
	}
	if err := c.ensureEnrolledStudentsAvailable(course); err != nil {
		return err
	}

	if err := c.repository.Update(course); err != nil {
		return err
	}

	if course.Credits != oldCredits {

		enrollments, err := c.enrollmentRepository.GetByCourse(
			course.Code,
		)

		if err != nil {
			return err
		}

		for _, enrollment := range enrollments {

			if err := recalculateStudentGPA(enrollment.StudentID, c.studentRepository, c.enrollmentRepository); err != nil {
				return err
			}
		}
	}

	return c.repository.Update(course)
}

func (c *courseService) ensureProfessorAvailable(course *models.Course) error {
	if course.ProfessorID == nil || *course.ProfessorID == "" {
		return nil
	}
	professorCourses, err := c.repository.GetByProfessor(*course.ProfessorID)
	if err != nil {
		return err
	}
	for i := range professorCourses {
		other := &professorCourses[i]
		if other.Code != course.Code && schedulesOverlap(course, other) {
			return fmt.Errorf("professor has a schedule conflict with %s", other.Code)
		}
	}
	return nil
}

func (c *courseService) ensureEnrolledStudentsAvailable(course *models.Course) error {
	roster, err := c.enrollmentRepository.GetByCourse(course.Code)
	if err != nil {
		return err
	}
	for _, rosterEntry := range roster {
		enrollments, err := c.enrollmentRepository.GetByStudent(rosterEntry.StudentID)
		if err != nil {
			return err
		}
		for _, enrollment := range enrollments {
			if enrollment.Course != nil && enrollment.CourseCode != course.Code && schedulesOverlap(course, enrollment.Course) {
				return fmt.Errorf("rescheduling conflicts with %s for student %d", enrollment.CourseCode, rosterEntry.StudentID)
			}
		}
	}
	return nil
}

func (c *courseService) Delete(code string) error {

	if _, err := c.repository.GetByCode(code); err != nil {
		return err
	}

	enrollments, err := c.enrollmentRepository.GetByCourse(code)

	if err != nil {
		return err
	}

	if err := c.enrollmentRepository.DeleteByCourse(code); err != nil {
		return err
	}

	if err := c.repository.Delete(code); err != nil {
		return err
	}

	for _, enrollment := range enrollments {

		if err := recalculateStudentGPA(enrollment.StudentID, c.studentRepository, c.enrollmentRepository); err != nil {
			return err
		}
	}

	return nil
}
func (c *courseService) Search(title string) ([]models.Course, error) {
	return c.repository.Search(title)
}

func (c *courseService) FilterByCredits(credits int) ([]models.Course, error) {
	return c.repository.FilterByCredits(credits)
}
