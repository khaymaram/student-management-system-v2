// Package routes defines the backend HTTP routing for the student API.
// It configures CORS, the /api prefix, and the /students route group.
package routes

import (
	"net/http"

	"sms-backend/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(
	studentHandler *handlers.StudentHandler,
	courseHandler *handlers.CourseHandler,
	enrollmentHandler *handlers.EnrollmentHandler,
) *gin.Engine {

	router := gin.Default()

	// Enable CORS so the browser frontend can talk to the backend during development.
	// Without this, modern browsers would block the API request as a cross-origin call.
	router.Use(cors.Default())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "pong",
			},
		)
	})

	// Group all student endpoints under /api for a clean backend URL structure.
	// The frontend calls /api/students, /api/students/:id, and related routes.
	api := router.Group("/api")
	{

		students := api.Group("/students")
		{

			students.GET(
				"",
				studentHandler.GetAll,
			)

			students.POST(
				"",
				studentHandler.Create,
			)

			students.GET(
				"/:studentId",
				studentHandler.GetByID,
			)

			students.PUT(
				"/:studentId",
				studentHandler.Update,
			)

			students.DELETE(
				"/:studentId",
				studentHandler.Delete,
			)

			students.GET(
				"/search",
				studentHandler.Search,
			)

			students.GET(
				"/grade/:grade",
				studentHandler.FilterByGrade,
			)

			students.GET(
				"/honors",
				studentHandler.GetHonors,
			)

			// Enrollment routes — a student's list of courses.
			students.GET(
				"/:studentId/enrollments",
				enrollmentHandler.GetByStudent,
			)

			students.POST(
				"/:studentId/enrollments",
				enrollmentHandler.Enroll,
			)

			students.PUT(
				"/:studentId/enrollments/:courseCode",
				enrollmentHandler.UpdateGrade,
			)

			students.DELETE(
				"/:studentId/enrollments/:courseCode",
				enrollmentHandler.Unenroll,
			)

		}

		courses := api.Group("/courses")
		{

			courses.GET(
				"",
				courseHandler.GetAll,
			)

			courses.POST(
				"",
				courseHandler.Create,
			)

			courses.GET(
				"/search",
				courseHandler.Search,
			)

			courses.GET(
				"/credits/:credits",
				courseHandler.FilterByCredits,
			)
			
			courses.GET(
				"/:code",
				courseHandler.GetByCode,
			)

			courses.PUT(
				"/:code",
				courseHandler.Update,
			)

			courses.DELETE(
				"/:code",
				courseHandler.Delete,
			)

			// Roster route — students enrolled in this course.
			courses.GET(
				"/:code/roster",
				enrollmentHandler.GetByCourse,
			)

		}

	}

	return router
}
