package handlers

import (
	"net/http"

	"sms-backend/internal/helpers"
	"sms-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type MajorHandler struct{ repository repositories.MajorRepository }

func NewMajorHandler(repository repositories.MajorRepository) *MajorHandler {
	return &MajorHandler{repository: repository}
}

func (h *MajorHandler) GetAll(c *gin.Context) {
	majors, err := h.repository.GetAll()
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.SuccessResponse(c, http.StatusOK, majors)
}
