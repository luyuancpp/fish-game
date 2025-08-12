package ws

import (
	"context"
	"fish-game/apps/user/user"
	"fish-game/ws/proto"
	"google.golang.org/protobuf/proto"
	"log"
)

func HandleShoot(c *Client, msg *ws.WSMessage) {
	var shoot ws.ShootRequest
	if err := proto.Unmarshal(msg.Data, &shoot); err != nil {
		log.Println("Failed to parse ShootRequest:", err)
		return
	}

	c.RoomHub.mu.Lock()
	defer c.RoomHub.mu.Unlock()

	for i, fish := range c.RoomHub.Fishes {
		if Hit(fish, shoot.X, shoot.Y) {
			// 删除鱼
			c.RoomHub.Fishes = append(c.RoomHub.Fishes[:i], c.RoomHub.Fishes[i+1:]...)

			// 广播 fish_killed
			killed := &ws.FishKilled{
				FishId:   fish.Id,
				ByUserId: AtoiSafe(c.UserID),
			}
			data, _ := proto.Marshal(killed)
			wrapped, _ := proto.Marshal(&ws.WSMessage{
				Event: "fish_killed",
				Data:  data,
			})
			c.RoomHub.Broadcast <- wrapped

			// 加金币
			go func(uid string) {
				reply, err := c.UserClient.AddGold(context.Background(), &user.AddGoldRequest{
					Uid:    uid,
					Amount: 10,
				})
				if err != nil {
					log.Println("AddGold error:", err)
					return
				}

				goldMsg := &ws.GoldUpdate{
					UserId: uid,
					Gold:   reply.Gold,
				}
				data, _ := proto.Marshal(goldMsg)
				wrapped, _ := proto.Marshal(&ws.WSMessage{
					Event: "gold_update",
					Data:  data,
				})
				c.RoomHub.Broadcast <- wrapped
			}(c.UserID)

			break
		}
	}
}
