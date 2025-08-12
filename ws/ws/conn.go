package ws

import (
	ws "fish-game/ws/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
)

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

		var wsMsg ws.WSMessage
		err = proto.Unmarshal(msg, &wsMsg)
		if err != nil {
			log.Println("Protobuf decode error:", err)
			continue
		}

		if handler, ok := handlers[wsMsg.Event]; ok {
			handler(c, &wsMsg)
		} else {
			log.Println("❓ 未知事件:", wsMsg.Event)
		}
	}

}

func (c *Client) writePump() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
