package websocket

import "github.com/i474232898/chatserver/internal/app/dto"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan dto.MessageDTO
	register   chan *Client
	unregister chan *Client
}

func (h Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case msg := <-h.broadcast:
			for cl := range h.clients {
				cl.Send <- msg
			}
		case client := <-h.unregister:
			delete(h.clients, client)
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan dto.MessageDTO),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
