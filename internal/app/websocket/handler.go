package websocket

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var hub = NewHub()

func init() {
	go hub.Run()
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer conn.Close()

	client := Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.register <- &client

	go client.Write()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message:", err)
			break
		}
		P(messageType, "<messageType")
		hub.broadcast <- message
	}
}
