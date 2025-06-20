package main

import (
	"fmt"
	"time"

	"github.com/i474232898/chatserver/configs"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/server"
)

func main() {
	cfg := configs.New()
	fmt.Println("Waiting for database to be ready...") //todo: add retry for db connection
	time.Sleep(5 * time.Second)

	_, err := repositories.GetPool(cfg)
	if err != nil {
		panic("Can't connect to db")
	}

	srvr := server.NewServer()
	srvr.Start(cfg.Port)
}
