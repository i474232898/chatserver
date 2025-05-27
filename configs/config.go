package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type dbConfig struct {
	User     string
	Password string
	Name     string
	Schema   string
	Host     string
	Port     string
}
type AppConfigs struct {
	Port        string
	DatabaseURL string
	DB          dbConfig
}

var config AppConfigs

func New() *AppConfigs {
	err := godotenv.Load()
	if err != nil {
		panic("incomplete environment vars")
	}
	config = AppConfigs{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DB: dbConfig{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_DB"),
			Schema:   os.Getenv("DB_SCHEMA"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
	}

	return &config
}
