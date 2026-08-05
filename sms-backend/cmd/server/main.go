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
	enrollmentRepository := repositories.NewEnrollmentRepository(db)
	studentRepository := repositories.NewStudentRepository(db)
	financeRepository := repositories.NewFinanceRepository(db)
	majorRepository := repositories.NewMajorRepository(db)

	studentService :=
		services.NewStudentService(studentRepository, enrollmentRepository, financeRepository, majorRepository)

	studentHandler :=
		handlers.NewStudentHandler(studentService)
	majorHandler := handlers.NewMajorHandler(majorRepository)

	financeService :=
		services.NewFinanceService(financeRepository)

	financeHandler :=
		handlers.NewFinanceHandler(financeService)

	courseRepository :=
		repositories.NewCourseRepository(db)

	courseService :=
		services.NewCourseService(courseRepository, enrollmentRepository, studentRepository)

	courseHandler :=
		handlers.NewCourseHandler(courseService)

	enrollmentService :=
		services.NewEnrollmentService(enrollmentRepository, studentRepository, courseRepository)

	enrollmentHandler :=
		handlers.NewEnrollmentHandler(enrollmentService)

	professorRepository :=
		repositories.NewProfessorRepository(db)

	professorService :=
		services.NewProfessorService(professorRepository, courseRepository)

	professorHandler :=
		handlers.NewProfessorHandler(professorService)

	router :=
		routes.Setup(
			studentHandler,
			courseHandler,
			enrollmentHandler,
			professorHandler,
			financeHandler,
			majorHandler,
		)

	err = router.Run(
		":" + cfg.Server.Port,
	)

	if err != nil {
		log.Fatal(err)
	}
}
