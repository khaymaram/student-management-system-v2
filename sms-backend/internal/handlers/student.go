// Package handlers implements HTTP handlers for student routes.
// It translates incoming requests into service calls and returns
// JSON responses using the helpers package.
package handlers

import (
	"net/http"
	"strconv"

	"sms-backend/internal/dto"
	"sms-backend/internal/helpers"
	"sms-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service services.StudentService
}

func NewStudentHandler(service services.StudentService) *StudentHandler {

	return &StudentHandler{
		service: service,
	}
}

// GET /students returns every student record.
func (h *StudentHandler) GetAll(c *gin.Context) {

	students, err := h.service.GetAll()

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
		students,
	)
}

// GET /students/:id
func (h *StudentHandler) GetByID(c *gin.Context) {

	studentId, err := strconv.Atoi(c.Param("studentId"))

	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid id",
		)
		return
	}

	student, err := h.service.GetByID(studentId)

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
		student,
	)
}

// POST /students creates a new student from the JSON body sent by the frontend.
// The handler reads the request data, passes it to the service layer, and returns
// a standardized success or error response back to the browser.
func (h *StudentHandler) Create(c *gin.Context) {

	var req dto.CreateStudentRequest

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
		"student created",
	)
}

// PUT /students/:id
func (h *StudentHandler) Update(c *gin.Context) {
	studentId, err := strconv.Atoi(c.Param("studentId"))
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid id",
		)
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	err = h.service.Update(studentId, req)
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
		"student updated",
	)
}

// DELETE /students/:id
func (h *StudentHandler) Delete(c *gin.Context) {
	studentId, err := strconv.Atoi(c.Param("studentId"))
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid id",
		)
		return
	}

	err = h.service.Delete(studentId)
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
		"student deleted",
	)
}

// GET /students/search?name=john
func (h *StudentHandler) Search(c *gin.Context) {
	name := c.Query("name")

	students, err := h.service.Search(name)
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
		students,
	)
}

// GET /students/grade/:grade
func (h *StudentHandler) FilterByGrade(c *gin.Context) {

	grade, err := strconv.Atoi(c.Param("grade"))

	if err != nil {

		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid grade",
		)

		return
	}

	students, err := h.service.FilterByGrade(grade)

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
		students,
	)
}

// GET /students/honors
func (h *StudentHandler) GetHonors(c *gin.Context) {

	students, err := h.service.GetHonors()

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
		students,
	)
}
