// Package handlers implements HTTP handlers for enrollment routes.
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

type EnrollmentHandler struct {
	service services.EnrollmentService
}

func NewEnrollmentHandler(service services.EnrollmentService) *EnrollmentHandler {

	return &EnrollmentHandler{
		service: service,
	}
}

// POST /students/:studentId/enrollments enrolls a student in a course.
func (h *EnrollmentHandler) Enroll(c *gin.Context) {

	studentId, err := strconv.Atoi(c.Param("studentId"))

	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var req dto.CreateEnrollmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Enroll(studentId, req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusCreated, "student enrolled")
}

// DELETE /students/:studentId/enrollments/:courseCode removes a student from a course.
func (h *EnrollmentHandler) Unenroll(c *gin.Context) {

	studentId, err := strconv.Atoi(c.Param("studentId"))

	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	courseCode := c.Param("courseCode")

	if err := h.service.Unenroll(studentId, courseCode); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "student unenrolled")
}

// PUT /students/:studentId/enrollments/:courseCode records/updates a final grade.
func (h *EnrollmentHandler) UpdateGrade(c *gin.Context) {

	studentId, err := strconv.Atoi(c.Param("studentId"))

	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	courseCode := c.Param("courseCode")

	var req dto.UpdateEnrollmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateGrade(studentId, courseCode, req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "grade updated")
}

// GET /enrollments lists every enrollment across the system.
func (h *EnrollmentHandler) GetAll(c *gin.Context) {
	enrollments, err := h.service.GetAll()

	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, enrollments)
}

// GET /students/:studentId/enrollments lists every course a student is taking.
func (h *EnrollmentHandler) GetByStudent(c *gin.Context) {

	studentId, err := strconv.Atoi(c.Param("studentId"))

	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	enrollments, err := h.service.GetByStudent(studentId)

	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, enrollments)
}

// GET /courses/:code/roster lists every student enrolled in a course.
func (h *EnrollmentHandler) GetByCourse(c *gin.Context) {

	courseCode := c.Param("code")

	enrollments, err := h.service.GetByCourse(courseCode)

	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, enrollments)
}
