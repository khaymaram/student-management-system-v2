// InitStudentHandler wires the database, repository, service, and handler
// layers together for student request handling.
package initializers

// import (
// 	"sms-backend/internal/handlers"
// 	"sms-backend/internal/repositories"
// 	"sms-backend/internal/services"

// 	"gorm.io/gorm"
// )

// InitStudentHandler wires the repository, service, and handler together.
// func InitStudentHandler(
// 	db *gorm.DB,
// ) *handlers.StudentHandler {

// 	studentRepository :=
// 		repositories.NewStudentRepository(db)

// 	studentService :=
// 		services.NewStudentService(
// 			studentRepository,
// 		)

// 	return handlers.NewStudentHandler(
// 		studentService,
// 	)
// }
