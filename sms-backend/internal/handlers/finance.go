package handlers

import (
	"net/http"
	"strconv"

	"sms-backend/internal/dto"
	"sms-backend/internal/helpers"
	"sms-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	service services.FinanceService
}

func NewFinanceHandler(service services.FinanceService) *FinanceHandler {
	return &FinanceHandler{
		service: service,
	}
}

func (h *FinanceHandler) GetAll(c *gin.Context) {
	finances, err := h.service.GetAllFinances()
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, finances)
}

func (h *FinanceHandler) Get(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("studentId"))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	finance, err := h.service.GetFinance(studentID)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, finance)
}

func (h *FinanceHandler) Update(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("studentId"))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var req dto.UpdateFinanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateFinance(studentID, req); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "finance updated")
}
