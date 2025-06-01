package websocket

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	pingWait = 5 * time.Second
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingWait)
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
		ticker.Stop()
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
		case <-ticker.C:
			err := c.Conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				c.Hub.unregister <- c
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
