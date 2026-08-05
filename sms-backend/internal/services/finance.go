package services

import (
	// "errors"

	"sms-backend/internal/dto"
	"sms-backend/internal/repositories"
)

type FinanceService interface {
	GetAllFinances() ([]dto.FinanceResponse, error)

	GetFinance(studentID int) (*dto.FinanceResponse, error)

	UpdateFinance(studentID int, req dto.UpdateFinanceRequest) error

	GetPaginated(page int, limit int)([]dto.FinanceResponse, int64, error)
}

const (
	InStateTuition  = 12000.0
	OutStateTuition = 25000.0
)

type financeService struct {
	repo repositories.FinanceRepository
}

func NewFinanceService(repo repositories.FinanceRepository) FinanceService {
	return &financeService{
		repo: repo,
	}
}

func (s *financeService) GetPaginated(
	page int,
	limit int,
) ([]dto.FinanceResponse, int64, error) {

	finances, total, err := s.repo.GetPaginated(page, limit)

	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.FinanceResponse, 0, len(finances))

	for _, finance := range finances {

		studentName := ""

		if finance.Student != nil {
			studentName = finance.Student.Name
		}

		remaining := finance.Tuition -
			finance.Scholarship -
			finance.Paid

		if remaining < 0 {
			remaining = 0
		}

		responses = append(responses, dto.FinanceResponse{
			StudentID:   finance.StudentID,
			StudentName: studentName,
			Tuition:     finance.Tuition,
			Scholarship: finance.Scholarship,
			Paid:        finance.Paid,
			IsInState:   finance.IsInState,
			Remaining:   remaining,
		})
	}

	return responses, total, nil
}

func (s *financeService) GetAllFinances() ([]dto.FinanceResponse, error) {
	finances, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.FinanceResponse, 0, len(finances))
	for _, finance := range finances {
		studentName := ""
		if finance.Student != nil {
			studentName = finance.Student.Name
		}

		remaining := finance.Tuition - finance.Scholarship - finance.Paid
		if remaining < 0 {
			remaining = 0
		}

		responses = append(responses, dto.FinanceResponse{
			StudentID:   finance.StudentID,
			StudentName: studentName,
			Tuition:     finance.Tuition,
			Scholarship: finance.Scholarship,
			Paid:        finance.Paid,
			IsInState:   finance.IsInState,
			Remaining:   remaining,
		})
	}

	return responses, nil
}

func (s *financeService) GetFinance(studentID int) (*dto.FinanceResponse, error) {

	finance, err := s.repo.GetByStudentID(studentID)

	if err != nil {
		return nil, err
	}

	remaining := finance.Tuition - finance.Scholarship - finance.Paid
	if remaining < 0 {
		remaining = 0
	}

	return &dto.FinanceResponse{
		StudentID: finance.StudentID,
		StudentName: func() string {
			if finance.Student != nil {
				return finance.Student.Name
			}
			return ""
		}(),
		Tuition:     finance.Tuition,
		Scholarship: finance.Scholarship,
		Paid:        finance.Paid,
		IsInState:   finance.IsInState,
		Remaining:   remaining,
	}, nil
}

func (s *financeService) UpdateFinance(studentID int, req dto.UpdateFinanceRequest) error {

	finance, err := s.repo.GetByStudentID(studentID)

	if err != nil {
		return err
	}

	finance.IsInState = req.IsInState
	finance.Scholarship = req.Scholarship
	finance.Paid = req.Paid

	if finance.IsInState {
		finance.Tuition = InStateTuition
	} else {
		finance.Tuition = OutStateTuition
	}

	return s.repo.Update(finance)
}
