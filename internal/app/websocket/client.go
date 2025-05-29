package websocket

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	P        = fmt.Println
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) Write() {
	defer func() {

		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{}) //send close frame
				return
			}
			if err := c.Conn.WriteMessage(1, msg); err != nil {
				slog.Debug(err.Error())
				return
			}
		}
	}
}

func (c *Client) Read() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message:", err)
			break
		}
		c.Hub.broadcast <- message
	}
}
