package main

import (
	"github.com/i474232898/chatserver/configs"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/server"
)

func main() {
	cfg := configs.New()
	_, err := repositories.GetPool(cfg)
	if err != nil {
		panic("Can't connect to db")
	}

	srvr := server.NewServer()
	srvr.Start(cfg.Port)
}
