package ws

import (
	"context"
	"fish-game/apps/user/user"
	"fish-game/ws/proto"
	"google.golang.org/protobuf/proto"
	"log"
	"math/rand"
	"time"
)

func HandleUseSkill(c *Client, msg *ws.WSMessage) {
	var req ws.UseSkillRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		log.Println("❌ UseSkillRequest 解析失败:", err)
		return
	}

	log.Printf("🧠 玩家 %s 使用技能: %s", req.UserId, req.SkillType)

	switch req.SkillType {

	case "freeze":
		c.RoomHub.SetFreeze(5 * time.Second)
		broadcastSkillUse("freeze", c.UserID, c)
		broadcastFreeze(5, c)

	case "missile":
		c.RoomHub.mu.Lock()
		killedFishes := c.RoomHub.Fishes
		c.RoomHub.Fishes = nil
		c.RoomHub.mu.Unlock()

		for _, fish := range killedFishes {
			killed := &ws.FishKilled{
				FishId:   fish.Id,
				ByUserId: AtoiSafe(c.UserID),
			}
			data, _ := proto.Marshal(killed)
			wrapped := &ws.WSMessage{
				Event: "fish_killed",
				Data:  data,
			}
			encoded, _ := proto.Marshal(wrapped)
			c.RoomHub.Broadcast <- encoded
		}

		broadcastSkillUse("missile", c.UserID, c)
		go rewardGold(c, int32(len(killedFishes)*10))

	case "lock_on":
		c.RoomHub.mu.Lock()
		if len(c.RoomHub.Fishes) == 0 {
			c.RoomHub.mu.Unlock()
			return
		}
		target := c.RoomHub.Fishes[rand.Intn(len(c.RoomHub.Fishes))]
		c.RoomHub.mu.Unlock()

		lockMsg := &ws.FishLocked{
			UserId: c.UserID,
			FishId: target.Id,
		}
		data, _ := proto.Marshal(lockMsg)
		wrapped := &ws.WSMessage{
			Event: "fish_locked",
			Data:  data,
		}
		encoded, _ := proto.Marshal(wrapped)
		c.RoomHub.Broadcast <- encoded

		broadcastSkillUse("lock_on", c.UserID, c)

	case "speed_up":
		msg := &ws.SpeedUp{
			UserId:   c.UserID,
			Duration: 5,
		}
		data, _ := proto.Marshal(msg)
		wrapped := &ws.WSMessage{
			Event: "speed_up",
			Data:  data,
		}
		encoded, _ := proto.Marshal(wrapped)
		c.RoomHub.Broadcast <- encoded

		broadcastSkillUse("speed_up", c.UserID, c)

	default:
		log.Printf("⚠️ 未知技能类型: %s", req.SkillType)
	}
}

func broadcastSkillUse(skillType, userID string, c *Client) {
	skillUsed := &ws.SkillUsed{
		UserId:    userID,
		SkillType: skillType,
	}
	data, _ := proto.Marshal(skillUsed)
	wrapped := &ws.WSMessage{
		Event: "skill_used",
		Data:  data,
	}
	encoded, _ := proto.Marshal(wrapped)
	c.RoomHub.Broadcast <- encoded
}

func broadcastFreeze(durationSec int, c *Client) {
	msg := &ws.FishFreeze{
		DurationMs: int32(durationSec),
	}
	data, _ := proto.Marshal(msg)
	wrapped := &ws.WSMessage{
		Event: "fish_freeze",
		Data:  data,
	}
	encoded, _ := proto.Marshal(wrapped)
	c.RoomHub.Broadcast <- encoded
}

func rewardGold(c *Client, amount int32) {
	reply, err := c.UserClient.AddGold(context.Background(), &user.AddGoldRequest{
		Uid:    c.UserID,
		Amount: amount,
	})
	if err != nil {
		log.Println("❌ AddGold error:", err)
		return
	}

	goldMsg := &ws.GoldUpdate{
		UserId: c.UserID,
		Gold:   reply.Gold,
	}
	data, _ := proto.Marshal(goldMsg)
	wrapped := &ws.WSMessage{
		Event: "gold_update",
		Data:  data,
	}
	encoded, _ := proto.Marshal(wrapped)
	c.RoomHub.Broadcast <- encoded
}
