package websocket

type Hub struct {
	clients   []*Client
	broadcast chan []byte
	register  chan *Client
}

func (h Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients = append(h.clients, client)

		case msg := <-h.broadcast:
			for _, cl := range h.clients {
				cl.Send <- msg
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		broadcast: make(chan []byte),
		register:  make(chan *Client),
		clients:   make([]*Client, 0),
	}
}
