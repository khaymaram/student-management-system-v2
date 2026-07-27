package config

import (
	"os"
	"strconv"
)

// Config holds the server and database settings loaded from environment variables.
type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func Load() *Config {

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}
}
