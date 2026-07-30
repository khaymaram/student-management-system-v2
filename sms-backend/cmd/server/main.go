// Command server is the backend entrypoint for the student management API.
// It loads configuration, connects to the database, constructs the
// repository/service/handler stack, and starts the Gin web server.
package main

import (
	"log"

	"sms-backend/internal/handlers"
	"sms-backend/internal/initializers"
	"sms-backend/internal/pkg/config"
	"sms-backend/internal/repositories"
	"sms-backend/internal/routes"
	"sms-backend/internal/services"
)

func main() {
	// Load environment variables, connect to the database, build the app stack, and start the server.

	initializers.LoadEnv()

	cfg := config.Load()

	db, err := initializers.ConnectDB(cfg)

	if err != nil {
		log.Fatal(err)
	}

	studentRepository :=
		repositories.NewStudentRepository(db)

	studentService :=
		services.NewStudentService(studentRepository)

	studentHandler :=
		handlers.NewStudentHandler(studentService)

	courseRepository :=
		repositories.NewCourseRepository(db)

	courseService :=
		services.NewCourseService(courseRepository)

	courseHandler :=
		handlers.NewCourseHandler(courseService)

	enrollmentRepository :=
		repositories.NewEnrollmentRepository(db)

	enrollmentService :=
		services.NewEnrollmentService(enrollmentRepository, studentRepository, courseRepository)

	enrollmentHandler :=
		handlers.NewEnrollmentHandler(enrollmentService)
		
	professorRepository :=
		repositories.NewProfessorRepository(db)

	professorService :=
		services.NewProfessorService(professorRepository)

	professorHandler :=
		handlers.NewProfessorHandler(professorService)

	router :=
		routes.Setup(
			studentHandler,
			courseHandler,
			enrollmentHandler,
			professorHandler,
		)

	err = router.Run(
		":" + cfg.Server.Port,
	)

	if err != nil {
		log.Fatal(err)
	}
}
