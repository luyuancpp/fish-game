package ws

import (
	"fish-game/ws/proto"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// HandleUseSkill 处理玩家使用技能的事件
func HandleUseSkill(c *Client, msg *ws.WSMessage) {
	var req ws.UseSkillRequest
	if err := proto.Unmarshal(msg.Data, &req); err != nil {
		log.Println("❌ UseSkillRequest 解析失败:", err)
		return
	}

	log.Printf("🧊 玩家 %s 使用技能: %s", req.UserId, req.SkillType)

	switch req.SkillType {
	case "freeze":
		// 冻结鱼的移动 N 秒（这里是 5 秒）
		duration := 5 * time.Second
		c.RoomHub.SetFreeze(duration)

		// 广播 skill_used
		skillUsed := &ws.SkillUsed{
			UserId:    req.UserId,
			SkillType: req.SkillType,
		}
		data1, _ := proto.Marshal(skillUsed)
		wrapped1 := &ws.WSMessage{
			Event: "skill_used",
			Data:  data1,
		}
		out1, _ := proto.Marshal(wrapped1)
		c.RoomHub.Broadcast <- out1

		// 广播 fish_freeze
		freezeMsg := &ws.FishFreeze{
			DurationMs: int32(duration / time.Second),
		}
		data2, _ := proto.Marshal(freezeMsg)
		wrapped2 := &ws.WSMessage{
			Event: "fish_freeze",
			Data:  data2,
		}
		out2, _ := proto.Marshal(wrapped2)
		c.RoomHub.Broadcast <- out2

	default:
		log.Println("⚠️ 未知技能类型:", req.SkillType)
	}
}
