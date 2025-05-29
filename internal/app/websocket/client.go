package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var P = fmt.Println

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c Client) Write() {
	for v := range c.Send {
		c.Conn.WriteMessage(1, v)
	}
}
