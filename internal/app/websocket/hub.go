package websocket

type Hub struct {
	clients  []*Client
	message  chan []byte
	register chan *Client
}

func (h Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients = append(h.clients, client)

		case msg := <-h.message:
			for _, cl := range h.clients {
				cl.Write(msg)
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		message:  make(chan []byte),
		register: make(chan *Client),
		clients:  make([]*Client, 0),
	}
}
