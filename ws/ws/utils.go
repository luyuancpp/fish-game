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

func broadcastCooldown(skillType, userID string, duration time.Duration, c *Client) {
	msg := &ws.SkillCooldown{
		UserId:      userID,
		SkillType:   skillType,
		DurationSec: int32(duration.Seconds()),
	}
	broadcastToAll("skill_cooldown", msg, c)
}

func hasItem(c *Client, itemID string) bool {
	// 假设有背包系统：
	return c.Inventory[itemID] > 0
}

func consumeItem(c *Client, itemID string) {
	if c.Inventory[itemID] > 0 {
		c.Inventory[itemID]--
	}
}

func sendTip(c *Client, content string) {
	msg := &ws.TipMessage{Text: content}
	broadcastToUser("tip", msg, c)
}

func getSkillLevel(c *Client, skillType string) int {
	if c.SkillLevels == nil {
		return 0
	}
	return c.SkillLevels[skillType]
}

func broadcastItemUpdate(c *Client, itemID string, count int) {
	msg := &ws.ItemUpdate{
		UserId: c.UserID,
		ItemId: itemID,
		Count:  int32(count),
	}
	broadcastToUser("item_update", msg, c)
}

// 向当前客户端单独发送一条消息（一般用于提示、私有反馈）
func broadcastToUser(event string, message proto.Message, c *Client) {
	data, err := proto.Marshal(message)
	if err != nil {
		log.Println("❌ Protobuf 编码失败:", err)
		return
	}

	wrapped := &ws.WSMessage{
		Event: event,
		Data:  data,
	}
	encoded, err := proto.Marshal(wrapped)
	if err != nil {
		log.Println("❌ 包装 WSMessage 编码失败:", err)
		return
	}

	// 发送给当前用户（非广播）
	select {
	case c.Send <- encoded:
	default:
		log.Println("⚠️ Send 队列满，消息丢弃")
	}
}

func broadcastToAllExcept(event string, message proto.Message, room *RoomHub, exceptUserID string) {
	data, err := proto.Marshal(message)
	if err != nil {
		log.Println("❌ Protobuf 编码失败:", err)
		return
	}

	wrapped := &ws.WSMessage{
		Event: event,
		Data:  data,
	}
	encoded, err := proto.Marshal(wrapped)
	if err != nil {
		log.Println("❌ WSMessage 编码失败:", err)
		return
	}

	for client := range room.Clients {
		if client.UserID == exceptUserID {
			continue
		}
		select {
		case client.Send <- encoded:
		default:
			log.Printf("⚠️ 用户 %s 队列已满，跳过发送", client.UserID)
		}
	}
}

func sendErrorTip(userID, text string, code int32, c *Client) {
	tip := &ws.TipMessage{
		UserId: userID,
		Text:   text,
		Code:   code,
	}
	broadcastToUser("tip", tip, c)
}

func broadcastToAll(event string, message proto.Message, c *Client) {
	data, err := proto.Marshal(message)
	if err != nil {
		log.Println("❌ Protobuf 编码失败:", err)
		return
	}

	wrapped := &ws.WSMessage{
		Event: event,
		Data:  data,
	}
	encoded, err := proto.Marshal(wrapped)
	if err != nil {
		log.Println("❌ 包装 WSMessage 编码失败:", err)
		return
	}

	c.RoomHub.Broadcast <- encoded
}
