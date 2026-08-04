package dto

type UpdateFinanceRequest struct {
	Scholarship float64 `json:"scholarship"`

	Paid float64 `json:"paid"`

	IsInState bool `json:"isInState"`
}

type FinanceResponse struct {
	StudentID int `json:"studentId"`

	StudentName string `json:"studentName,omitempty"`

	Tuition float64 `json:"tuition"`

	Scholarship float64 `json:"scholarship"`

	Paid float64 `json:"paid"`

	IsInState bool `json:"isInState"`

	Remaining float64 `json:"remaining"`
}
