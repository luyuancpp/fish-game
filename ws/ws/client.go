package ws

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	UserID  string
	RoomID  string
	Conn    *websocket.Conn
	Send    chan []byte
	RoomHub *RoomHub
}

func (c *Client) readPump() {
	defer func() {
		c.RoomHub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Println("📥 Client sent:", string(msg))
		c.RoomHub.Broadcast <- msg
	}
}

func (c *Client) writePump() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
		log.Println("📤 Server broadcasting to client:", string(msg))
	}
}
