package ws

import (
	ws "fish-game/ws/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
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

		var wsMsg ws.WSMessage
		err = proto.Unmarshal(msg, &wsMsg)
		if err != nil {
			log.Println("Protobuf decode error:", err)
			continue
		}

		switch wsMsg.Event {
		case "shoot":
			var shoot ws.ShootRequest
			if err := proto.Unmarshal(wsMsg.Data, &shoot); err == nil {
				log.Printf("💥 用户 %s 发射子弹：ID=%d X=%.2f Y=%.2f\n", c.UserID, shoot.BulletId, shoot.X, shoot.Y)
				// 可以进行碰撞检测或广播
			}
		default:
			log.Println("❓ 未知事件：", wsMsg.Event)
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
