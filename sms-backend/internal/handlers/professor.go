package handlers

import (
	"net/http"

	"sms-backend/internal/dto"
	"sms-backend/internal/helpers"
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
