package main

import (
	"github.com/i474232898/chatserver/configs"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/server"
)

func main() {
	cnf := configs.New()
	err := repositories.Connect(cnf)
	if err != nil {
		panic("Can't connect to db")
	}

	srvr := server.NewServer()
	srvr.Start(cnf.Port)
}
