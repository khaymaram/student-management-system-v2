package handlers

import (
	"net/http"
	"strconv"

	"sms-backend/internal/dto"
	"sms-backend/internal/helpers"
	"sms-backend/internal/models"
	"sms-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ProfessorHandler struct {
	service services.ProfessorService
}

func NewProfessorHandler(
	service services.ProfessorService,
) *ProfessorHandler {

	return &ProfessorHandler{
		service: service,
	}
}

// GET /professors returns every student record.
func (h *ProfessorHandler) GetAll(c *gin.Context) {

	pageParam := c.Query("page")
	limitParam := c.Query("limit")

	if pageParam == "" && limitParam == "" {
		professors, err := h.service.GetAll()

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
			professors,
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

	professorId := c.Query("professorId")
	name := c.Query("name")

	professors, total, err := h.service.GetPaginated(
		page, limit, professorId, name,
	)

	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	response := dto.PaginationResponse[models.Professor]{
		Data:       professors,
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

// GET /professors/:code
func (h *ProfessorHandler) GetByID(c *gin.Context) {

	professorCode := c.Param("id")

	professor, err := h.service.GetByID(professorCode)

	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		professor,
	)
}

// POST /professors creates a new student from the JSON body sent by the frontend.
// The handler reads the request data, passes it to the service layer, and returns
// a standardized success or error response back to the browser.
func (h *ProfessorHandler) Create(c *gin.Context) {

	var req dto.CreateProfessorRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	err := h.service.Create(req)

	if err != nil {

		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	helpers.SuccessResponse(
		c,
		http.StatusCreated,
		"professor created",
	)
}

// PUT /professors/:id
func (h *ProfessorHandler) Update(c *gin.Context) {
	professorCode := c.Param("id")

	var req dto.UpdateProfessorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	err := h.service.Update(professorCode, req)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		"professor updated",
	)
}

// DELETE /professors/:code
func (h *ProfessorHandler) Delete(c *gin.Context) {
	professorCode := c.Param("id")

	err := h.service.Delete(professorCode)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		"professor deleted",
	)
}

// GET /professors/search?name=john
func (h *ProfessorHandler) Search(c *gin.Context) {
	name := c.Query("name")

	professors, err := h.service.Search(name)
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
		professors,
	)
}
