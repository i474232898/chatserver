package websocket

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/i474232898/chatserver/internal/app/common"
	"github.com/i474232898/chatserver/internal/app/services"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var hm = NewHubManager()

type WebsocketHandler struct {
	service services.ChatRoomService
}

func NewWebsocketHandler(serv services.ChatRoomService) *WebsocketHandler {
	return &WebsocketHandler{service: serv}
}

func (h *WebsocketHandler) JoinChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	roomId, err := strconv.Atoi(chi.URLParam(r, "roomID"))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	claims, err := common.ParseJWT(token, []byte("secret"))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	invited := h.service.IsUserInRoom(r.Context(), uint64(claims.ID), uint64(roomId))
	if !invited {
		slog.Info("Not invited")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	hub := hm.GetHub(roomId)
	client := Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.register <- &client

	go client.Read()
	go client.Write()
}
