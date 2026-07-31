package handlers

import (
	"net/http"
	"strconv"

	"sms-backend/internal/dto"
	"sms-backend/internal/helpers"
	"sms-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service services.CourseService
}

func NewCourseHandler(service services.CourseService) *CourseHandler {

	return &CourseHandler{
		service: service,
	}
}

// GET /courses returns every student record.
// func (h *CourseHandler) GetAll(c *gin.Context) {

// 	courses, err := h.service.GetAll()

// 	if err != nil {
// 		helpers.ErrorResponse(
// 			c,
// 			http.StatusInternalServerError,
// 			err.Error(),
// 		)
// 		return
// 	}

// 	helpers.SuccessResponse(
// 		c,
// 		http.StatusOK,
// 		courses,
// 	)
// }
func (h *CourseHandler) GetAll(c *gin.Context) {

    professorId := c.Query("professorId")

    if professorId != "" {

        courses, err := h.service.GetByProfessor(professorId)

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
            courses,
        )
        return
    }

    courses, err := h.service.GetAll()

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
        courses,
    )
}
func (h *CourseHandler) GetByProfessor(c *gin.Context) {

	professorId := c.Param("professorId")

	courses, err := h.service.GetByProfessor(professorId)

	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, courses)
}

// GET /courses/:code
func (h *CourseHandler) GetByCode(c *gin.Context) {

	courseCode := c.Param("code")

	course, err := h.service.GetByCode(courseCode)

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
		course,
	)
}

// POST /courses creates a new student from the JSON body sent by the frontend.
// The handler reads the request data, passes it to the service layer, and returns
// a standardized success or error response back to the browser.
func (h *CourseHandler) Create(c *gin.Context) {

	var req dto.CreateCourseRequest

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
		"course created",
	)
}

// PUT /courses/:id
func (h *CourseHandler) Update(c *gin.Context) {
	// courseId, err := strconv.Atoi(c.Param("courseId"))
	courseCode := c.Param("code")
	// if err != nil {
	// 	helpers.ErrorResponse(
	// 		c,
	// 		http.StatusBadRequest,
	// 		"invalid id",
	// 	)
	// 	return
	// }

	var req dto.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	err := h.service.Update(courseCode, req)
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
		"course updated",
	)
}

// DELETE /courses/:code
func (h *CourseHandler) Delete(c *gin.Context) {
	courseCode := c.Param("code")

	

	err := h.service.Delete(courseCode)
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
		"course deleted",
	)
}

// GET /courses/search?name=john
func (h *CourseHandler) Search(c *gin.Context) {
	title := c.Query("title")

	courses, err := h.service.Search(title)
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
		courses,
	)
}

// GET /courses/credits/:credits
func (h *CourseHandler) FilterByCredits(c *gin.Context) {

	credits, err := strconv.Atoi(c.Param("credits"))

	if err != nil {

		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid number of credits",
		)

		return
	}

	courses, err := h.service.FilterByCredits(credits)

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
		courses,
	)
}

