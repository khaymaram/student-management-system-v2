// Package initializers provides startup setup helpers for the backend.
// ConnectDB opens a GORM connection to the configured SQL Server database.
package initializers

import (
	"sms-backend/internal/pkg/config"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// ConnectDB opens the GORM connection to the SQL Server database.
// The database configuration is loaded from config values and used to form the DSN.
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {

	dsn := "server=localhost;database=test; integrated security=true"

	return gorm.Open(
		sqlserver.Open(dsn),
		&gorm.Config{},
	)
}
