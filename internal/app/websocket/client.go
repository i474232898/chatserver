package websocket

import (
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/services"
)

var (
	pongWait = 10 * time.Second
	pingWait = 5 * time.Second
)

type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	Send        chan dto.MessageDTO
	RoomId      uint64
	UserId      uint64
	RoomService services.ChatRoomService
}

//send messages to client
func (c *Client) Write(lastSentMessageId uint64) {
	ticker := time.NewTicker(pingWait)
	defer func() {
		c.Hub.unregister <- c
		err := c.Conn.Close()
		if err != nil {
			slog.Error("Error closing connection: " + err.Error())
		}
		ticker.Stop()
	}()

	//send all messages that client hasn't seen yet
	msgs, err := c.RoomService.GetMessages(context.Background(), c.RoomId, uint64(lastSentMessageId))
	if err != nil {
		slog.Error("Error getting messages: " + err.Error())
	}
	for _, msg := range msgs {
		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		if err != nil {
			slog.Error("Error writing message: " + err.Error())
			return
		}
	}

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}) //send close frame
				if err != nil {
					slog.Error("Error writing close message: " + err.Error())
				}
				return
			}
			if err := c.Conn.WriteMessage(1, []byte(msg.Content)); err != nil {
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

//receive messages from client
func (c *Client) Read() {
	defer func() {
		c.Hub.unregister <- c
		err := c.Conn.Close()
		if err != nil {
			slog.Error("Error closing connection: " + err.Error())
		}
	}()

	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		slog.Error("Error setting read deadline: " + err.Error())
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			slog.Error("Error setting read deadline: " + err.Error())
		}
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message: " + err.Error())
			break
		}
		//todo: should new context be created for every message save?
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		msgDto, err := c.RoomService.SaveMessage(ctx, c.RoomId, c.UserId, string(message))
		cancel()
		if err != nil {
			slog.Error("Error saving message: " + err.Error())
			continue
		}

		c.Hub.broadcast <- *msgDto
	}
}
