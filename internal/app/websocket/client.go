package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
)

var P = fmt.Println

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c Client) Write() {
	// for v := range c.Send {
	// 	c.Conn.WriteMessage(1, v)
	// }
	defer func() {
		c.Hub.unregister <- &c
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
