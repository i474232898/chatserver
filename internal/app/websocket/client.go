package websocket

import (
	"context"
	// "fmt"
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

func (c *Client) Write() {
	ticker := time.NewTicker(pingWait)
	lastSentMessageId, err := c.RoomService.GetUserRoomOffset(context.Background(), c.RoomId, c.UserId)
	if err != nil {
		slog.Error("Error getting user room offset", "error", err.Error())
	}
	// fmt.Println(lastSentMessageId, c.RoomId, c.UserId, "<<lastSentMessageId")
	defer func() {
		c.Hub.unregister <- c
		err := c.Conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
		ticker.Stop()

		err = c.RoomService.UpdateUserRoomOffset(context.Background(), c.RoomId, c.UserId, uint64(lastSentMessageId))
		if err != nil {
			slog.Error("Error updating user room offset", "error", err.Error())
		}
	}()

	//send all messages that client hasn't seen yet
	msgs, err := c.RoomService.GetMessages(context.Background(), c.RoomId, uint64(lastSentMessageId))
	if err != nil {
		slog.Error("Error getting messages", "error", err.Error())
	}
	for _, msg := range msgs {
		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		if err != nil {
			slog.Error("Error writing message", "error", err.Error())
			return
		}
		lastSentMessageId = uint64(msg.ID)
		// fmt.Println(lastSentMessageId, "222")
	}

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
			if err := c.Conn.WriteMessage(1, []byte(msg.Content)); err != nil {
				slog.Debug(err.Error())
				return
			}
			lastSentMessageId = uint64(msg.ID)
			// fmt.Println(lastSentMessageId, "<<<")
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
		//todo: should new context be created for every message save?
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		msgDto, err := c.RoomService.SaveMessage(ctx, c.RoomId, c.UserId, string(message))
		cancel()
		if err != nil {
			slog.Error("Error saving message", "error", err.Error())
			continue
		}

		c.Hub.broadcast <- *msgDto
	}
}
