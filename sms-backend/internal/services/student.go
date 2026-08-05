// Package services contains application business logic for students.
// It validates payloads, prevents duplicate student IDs, and orchestrates
// calls to the repository layer.
package services

import (
	"errors"

	"sms-backend/internal/dto"
	"sms-backend/internal/models"
	"sms-backend/internal/repositories"

	"gorm.io/gorm"
)

type StudentService interface {
	GetAll() ([]models.Student, error)

	GetPaginated(
		page int,
		limit int,
		grade *int, 
		honors bool,
		studentID *int,
		name string,
	) ([]models.Student, int64, error)

	GetByID(id int) (*models.Student, error)

	Create(dto.CreateStudentRequest) error

	Update(id int, req dto.UpdateStudentRequest) error

	Delete(id int) error

	Search(name string) ([]models.Student, error)

	FilterByGrade(grade int) ([]models.Student, error)

	GetHonors() ([]models.Student, error)
}

type studentService struct {
	repository           repositories.StudentRepository
	enrollmentRepository repositories.EnrollmentRepository
	financeRepository    repositories.FinanceRepository
}

func NewStudentService(repo repositories.StudentRepository, eRepo repositories.EnrollmentRepository, fRepo repositories.FinanceRepository) StudentService {

	return &studentService{
		repository:           repo,
		enrollmentRepository: eRepo,
		financeRepository:    fRepo,
	}
}

func (s *studentService) GetAll() ([]models.Student, error) {
	return s.repository.GetAll()
}

func (s *studentService) GetPaginated(
	page int,
	limit int,
	grade *int,
	honors bool,
	studentID *int,
	name string,
) ([]models.Student, int64, error) {
	return s.repository.GetPaginated(
		page, limit, grade, honors, studentID, name,
	)
}

func (s *studentService) GetByID(id int) (*models.Student, error) {
	return s.repository.GetByID(id)
}

func (s *studentService) Create(req dto.CreateStudentRequest) error {
	// Validate the request before saving it to the database.
	// This is the business-rule layer, where duplicate IDs and invalid GPA values are rejected.

	_, err := s.repository.GetByID(req.StudentID)

	if err == nil {
		return errors.New("student id already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if req.GPA < 0 || req.GPA > 4 {
		return errors.New("invalid GPA")
	}

	student := models.Student{
		ID:    req.StudentID,
		Name:  req.Name,
		Grade: req.Grade,
		GPA:   req.GPA,
	}

	if err := s.repository.Create(&student); err != nil {
		return err
	}

	finance := models.Finance{
		StudentID:   student.ID,
		Scholarship: req.Scholarship,
		Paid:        0,
		IsInState:   req.IsInState,
	}

	if finance.IsInState {
		finance.Tuition = InStateTuition
	} else {
		finance.Tuition = OutStateTuition
	}

	return s.financeRepository.Create(&finance)
}

func (s *studentService) Update(id int, req dto.UpdateStudentRequest) error {

	student, err := s.repository.GetByID(id)

	if err != nil {
		return err
	}

	if req.Name != "" {
		student.Name = req.Name
	}

	if req.Grade != 0 {
		student.Grade = int64(req.Grade)
	}

	if req.GPA != 0 {
		student.GPA = req.GPA
	}

	return s.repository.Update(student)
}

func (s *studentService) Delete(id int) error {
	if _, err := s.repository.GetByID(id); err != nil {
		return err
	}
	if err := s.enrollmentRepository.DeleteByStudent(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}

func (s *studentService) Search(name string) ([]models.Student, error) {
	return s.repository.Search(name)
}

func (s *studentService) FilterByGrade(grade int) ([]models.Student, error) {
	return s.repository.FilterByGrade(grade)
}

func (s *studentService) GetHonors() ([]models.Student, error) {
	return s.repository.GetHonors()
}
