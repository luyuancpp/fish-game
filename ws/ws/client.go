package ws

import (
	"context"
	"fish-game/apps/room/room"
	"fish-game/apps/user/user"
	ws "fish-game/ws/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
)

type Client struct {
	UserID     string
	RoomID     string
	Conn       *websocket.Conn
	Send       chan []byte
	RoomHub    *RoomHub
	RoomClient room.RoomClient
	UserClient user.UserClient // ✅ 加这一
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
			if err := proto.Unmarshal(wsMsg.Data, &shoot); err != nil {
				log.Println("Failed to parse ShootRequest:", err)
				continue
			}

			c.RoomHub.mu.Lock()
			for i, fish := range c.RoomHub.Fishes {
				if hit(fish, shoot.X, shoot.Y) {
					// ✅ 删除鱼
					c.RoomHub.Fishes = append(c.RoomHub.Fishes[:i], c.RoomHub.Fishes[i+1:]...)

					// ✅ 广播 fish_killed
					killed := &ws.FishKilled{
						FishId:   fish.Id,
						ByUserId: atoiSafe(c.UserID),
					}
					data, _ := proto.Marshal(killed)
					wrapper := &ws.WSMessage{
						Event: "fish_killed",
						Data:  data,
					}
					encoded, _ := proto.Marshal(wrapper)
					c.RoomHub.Broadcast <- encoded

					// ✅ 加金币（放在这）
					go func(uid string) {
						reply, err := c.UserClient.AddGold(context.Background(), &user.AddGoldRequest{
							Uid:    uid,
							Amount: 10, // 每条鱼奖励 10 金币
						})
						if err != nil {
							log.Println("❌ AddGold error:", err)
							return
						}

						log.Printf("💰 用户 %s 获得金币，当前余额：%d", uid, reply.Gold)

						goldMsg := &ws.GoldUpdate{
							UserId: uid,
							Gold:   reply.Gold,
						}
						data, _ := proto.Marshal(goldMsg)
						wrapped := &ws.WSMessage{
							Event: "gold_update",
							Data:  data,
						}
						encoded, _ := proto.Marshal(wrapped)
						c.RoomHub.Broadcast <- encoded
					}(c.UserID)

					break
				}
			}
			c.RoomHub.mu.Unlock()

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
