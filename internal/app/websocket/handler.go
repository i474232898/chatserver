package websocket

import (
	"log/slog"
	"net/http"

	// "github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	// "github.com/i474232898/chatserver/internal/app/common"
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

	client := Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.register <- &client

	go client.Read()
	go client.Write()
}

// ws/room/{roomID}?token=JWT
func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// token := r.URL.Query().Get("token")
	// roomId := chi.URLParam(r, "roomID")

	// claims, err := common.ParseJWT(token, []byte("secret"))
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	return
	// }

	// //check is user in room
	// repo	
}
