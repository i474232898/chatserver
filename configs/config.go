package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfigs struct {
	Port        string
	DatabaseURL string
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
	}

	return &config
}
