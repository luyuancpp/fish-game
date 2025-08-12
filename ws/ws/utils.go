package ws

import (
	"context"
	"fish-game/apps/user/user"
	ws "fish-game/ws/proto"
	"google.golang.org/protobuf/proto"
	"log"
	"math"
	"strconv"
	"time"
)

func AtoiSafe(s string) int32 {
	i, _ := strconv.Atoi(s)
	return int32(i)
}

func distance(x1, y1, x2, y2 float32) float32 {
	dx := x1 - x2
	dy := y1 - y2
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func Hit(fish *Fish, bulletX, bulletY float32) bool {
	dx := fish.X - bulletX
	dy := fish.Y - bulletY
	distance := float32(math.Hypot(float64(dx), float64(dy)))
	return distance <= fish.HitRadius
}

func broadcastSkillUse(skillType, userID string, c *Client) {
	skillUsed := &ws.SkillUsed{
		UserId:    userID,
		SkillType: skillType,
	}
	broadcastToAll("skill_used", skillUsed, c)
}

func broadcastToAll(event string, msg proto.Message, c *Client) {
	data, _ := proto.Marshal(msg)
	wrapped := &ws.WSMessage{
		Event: event,
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
	broadcastToAll("gold_update", goldMsg, c)
}

func (c *Client) IsSkillCoolingDown(skill string) bool {
	next := c.Cooldowns[skill]
	return time.Now().Before(next)
}

func (c *Client) SetSkillCooldown(skill string, duration time.Duration) {
	c.Cooldowns[skill] = time.Now().Add(duration)
}

func isCoolingDown(c *Client, skill string) bool {
	next, ok := c.Cooldowns[skill]
	return ok && time.Now().Before(next)
}

func setCooldown(c *Client, skill string, duration time.Duration) {
	c.Cooldowns[skill] = time.Now().Add(duration)
}
