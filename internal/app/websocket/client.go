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
		err := c.Conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}) //send close frame
				if err != nil {
					slog.Error("Error writing close message", "error", err)
				}
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
		err := c.Conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
	}()

	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		slog.Error("Error setting read deadline", "error", err.Error())
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			slog.Error("Error setting read deadline", "error", err.Error())
		}
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", "error", err.Error())
			break
		}
		c.Hub.broadcast <- message
	}
}
