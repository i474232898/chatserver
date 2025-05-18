package main

import (
	"os"

	"github.com/i474232898/chatserver/internal/app/server"
	"github.com/i474232898/chatserver/pkg/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("incomplete environment vars")
	}

	err = db.Connect()
	if err != nil {
		panic("Can't connect to db")
	}

	server.Start(os.Getenv("PORT"))
}
