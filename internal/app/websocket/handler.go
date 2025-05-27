package websocket

import (
	// "fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer conn.Close()

	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message:", err)
			break
		}

		// Append "abc11" to the received message
		response := string(message) + "abc11"

		// Send the modified message back to the client
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			slog.Error("Error writing message:", err)
			break
		}
	}

	// message := struct {
	// 	A string
	// }{}

	// err = conn.ReadJSON(&message)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	return
	// }
	// fmt.Println(message, "<<<")
}
