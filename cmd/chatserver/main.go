package main

import (
	"fmt"
	"os"

	"github.com/i474232898/chatserver/internal/app/server"
	"github.com/i474232898/chatserver/pkg/db"
	"github.com/joho/godotenv"
)

var p = fmt.Println

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("incomplete environment vars")
	}

	db.Connect()
	server.Start(os.Getenv("PORT"))
}
