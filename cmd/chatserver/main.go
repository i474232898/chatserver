package main

import (
	"os"

	"github.com/i474232898/chatserver/docs"
	"github.com/i474232898/chatserver/internal/app/server"
	"github.com/i474232898/chatserver/pkg/db"
	"github.com/joho/godotenv"
)

// @title           Chat Server API
// @version         1.0
// @description     A chat server API built with Go
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("incomplete environment vars")
	}

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Chat Server API"
	docs.SwaggerInfo.Description = "A chat server API built with Go"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	err = db.Connect()
	if err != nil {
		panic("Can't connect to db")
	}

	server.Start(os.Getenv("PORT"))
}
