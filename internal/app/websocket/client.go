package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var P = fmt.Println

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
}

func NewClient(h *Hub, conn *websocket.Conn) *Client {
	return &Client{
		Hub:  h,
		Conn: conn,
	}
}

func (c Client) Write(msg []byte) {
	// str := string(msg)
	// P(str, "<<write<")
	c.Conn.WriteMessage(1, msg)
}
