package main

import (
	"github.com/i474232898/chatserver/configs"
	"github.com/i474232898/chatserver/internal/app/server"
	"github.com/i474232898/chatserver/pkg/db"
)

func main() {
	cnf := configs.New()
	err := db.Connect(cnf)
	if err != nil {
		panic("Can't connect to db")
	}

	srvr := server.NewServer()
	srvr.Start(cnf.Port)
}
