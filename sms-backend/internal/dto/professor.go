package dto

type CreateProfessorRequest struct {
	ID string `json:"id" binding:"required"`

	Name string `json:"name" binding:"required"`
}

type UpdateProfessorRequest struct {
	Name string `json:"name" binding:"required"`
}