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
	pageParam := c.Query("page")
	limitParam := c.Query("limit")

	// Keep the original non-paginated endpoint behavior.
	if pageParam == "" && limitParam == "" {
		finances, err := h.service.GetAllFinances()

		if err != nil {
			helpers.ErrorResponse(
				c,
				http.StatusInternalServerError,
				err.Error(),
			)
			return
		}

		helpers.SuccessResponse(
			c,
			http.StatusOK,
			finances,
		)
		return
	}

	page := 1
	limit := 5

	if pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)

		if err != nil || parsedPage < 1 {
			helpers.ErrorResponse(
				c,
				http.StatusBadRequest,
				"page must be a positive integer",
			)
			return
		}

		page = parsedPage
	}

	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)

		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			helpers.ErrorResponse(
				c,
				http.StatusBadRequest,
				"limit must be between 1 and 100",
			)
			return
		}

		limit = parsedLimit
	}

	finances, total, err := h.service.GetPaginated(page, limit)

	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	response := dto.PaginationResponse[dto.FinanceResponse]{
		Data:       finances,
		Page:       page,
		PageSize:   limit,
		Total:      total,
		TotalPages: totalPages,
	}

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		response,
	)
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
