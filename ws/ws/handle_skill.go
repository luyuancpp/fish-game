package ws

import (
	"fish-game/config"
	"fish-game/ws/proto"
	"fmt"
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

	cf, ok := config.SkillConfigs[req.SkillType]
	if !ok {
		log.Printf("⚠️ 未知技能类型: %s", req.SkillType)
		return
	}

	// 检查道具数量
	if !hasItem(c, cf.RequiredItem) {
		sendTip(c, "❌ 道具不足，无法使用该技能")
		return
	}

	// 检查冷却
	if c.Cooldowns == nil {
		c.Cooldowns = make(map[string]time.Time)
	}
	if isCoolingDown(c, req.SkillType) {
		sendTip(c, fmt.Sprintf("⏳ 技能 [%s] 冷却中", req.SkillType))
		return
	}

	// 设置冷却 & 消耗道具
	setCooldown(c, req.SkillType, cf.Cooldown)
	consumeItem(c, cf.RequiredItem)
	broadcastCooldown(req.SkillType, req.UserId, cf.Cooldown, c)

	// 技能等级效果增强（可以从 c.SkillLevels[req.SkillType] 获取等级）
	//level := getSkillLevel(c, req.SkillType)

	switch req.SkillType {
	case "freeze":
		handleFreezeSkill(c, &req)
	case "missile":
		handleMissileSkill(c, &req)
	case "lock_on":
		handleLockOnSkill(c, &req)
	case "speed_up":
		handleSpeedUpSkill(c, &req)
	case "radar":
		handleRadarSkill(c, &req)
	case "slow_area":
		handleSlowAreaSkill(c, &req)
	case "magnet":
		handleMagnetSkill(c, &req, 150)
	case "emp_blast":
		handleEMPBlastSkill(c, &req, 150, 3*time.Second)
	case "invisible":
		handleInvisibleSkill(c, &req, 5*time.Second)
	case "grow_fish":
		handleGrowFishSkill(c, &req, 2.0)
	default:
		log.Printf("⚠️ 未知技能类型: %s", req.SkillType)
	}
}

func handleFreezeSkill(c *Client, req *ws.UseSkillRequest) {
	duration := 5 * time.Second
	c.RoomHub.SetFreeze(duration)

	broadcastSkillUse("freeze", req.UserId, c)

	msg := &ws.FishFreeze{
		DurationMs: int32(duration.Milliseconds()),
	}
	broadcastToAll("fish_freeze", msg, c)
}

func handleMissileSkill(c *Client, req *ws.UseSkillRequest) {
	c.RoomHub.mu.Lock()
	killedFishes := make([]*Fish, len(c.RoomHub.Fishes))
	copy(killedFishes, c.RoomHub.Fishes)
	c.RoomHub.Fishes = nil
	c.RoomHub.mu.Unlock()

	for _, fish := range killedFishes {
		killed := &ws.FishKilled{
			FishId:   fish.Id,
			ByUserId: AtoiSafe(req.UserId),
		}
		broadcastToAll("fish_killed", killed, c)
	}

	broadcastSkillUse("missile", req.UserId, c)
	go rewardGold(c, int32(len(killedFishes)*10))
}

func handleLockOnSkill(c *Client, req *ws.UseSkillRequest) {
	c.RoomHub.mu.Lock()
	defer c.RoomHub.mu.Unlock()

	if len(c.RoomHub.Fishes) == 0 {
		return
	}
	target := c.RoomHub.Fishes[rand.Intn(len(c.RoomHub.Fishes))]

	lockMsg := &ws.FishLocked{
		UserId: req.UserId,
		FishId: target.Id,
	}
	broadcastToAll("fish_locked", lockMsg, c)

	broadcastSkillUse("lock_on", req.UserId, c)
}

func handleSpeedUpSkill(c *Client, req *ws.UseSkillRequest) {
	duration := 5

	msg := &ws.SpeedUp{
		UserId:   req.UserId,
		Duration: int32(duration),
	}
	broadcastToAll("speed_up", msg, c)

	broadcastSkillUse("speed_up", req.UserId, c)
}

func handleRadarSkill(c *Client, req *ws.UseSkillRequest) {
	centerX, centerY := req.X, req.Y
	radius := float32(100)

	c.RoomHub.mu.Lock()
	defer c.RoomHub.mu.Unlock()

	var result []*ws.FishPosition
	for _, fish := range c.RoomHub.Fishes {
		if distance(fish.X, fish.Y, centerX, centerY) <= radius {
			result = append(result, &ws.FishPosition{
				Id: fish.Id,
				X:  fish.X,
				Y:  fish.Y,
			})
		}
	}

	resp := &ws.RadarResult{
		UserId:    req.UserId,
		CenterX:   centerX,
		CenterY:   centerY,
		Radius:    radius,
		FishFound: result,
	}
	broadcastToAll("radar_result", resp, c)

	broadcastSkillUse("radar", req.UserId, c)
}

func handleSlowAreaSkill(c *Client, req *ws.UseSkillRequest) {
	centerX, centerY := req.X, req.Y
	radius := float32(120)
	duration := 5 * time.Second

	c.RoomHub.mu.Lock()
	for _, fish := range c.RoomHub.Fishes {
		if distance(fish.X, fish.Y, centerX, centerY) <= radius {
			fish.SpeedX *= 0.5
			fish.SpeedY *= 0.5
		}
	}
	c.RoomHub.mu.Unlock()

	go func() {
		time.Sleep(duration)
		c.RoomHub.mu.Lock()
		defer c.RoomHub.mu.Unlock()

		for _, fish := range c.RoomHub.Fishes {
			if distance(fish.X, fish.Y, centerX, centerY) <= radius {
				fish.SpeedX *= 2
				fish.SpeedY *= 2
			}
		}
	}()

	msg := &ws.AreaSlowDown{
		CenterX:  centerX,
		CenterY:  centerY,
		Radius:   radius,
		Duration: int32(duration.Seconds()),
	}
	broadcastToAll("area_slow_down", msg, c)

	broadcastSkillUse("slow_area", req.UserId, c)
}

func handleMagnetSkill(c *Client, req *ws.UseSkillRequest, radius float32) {
	c.RoomHub.mu.Lock()
	defer c.RoomHub.mu.Unlock()

	for _, fish := range c.RoomHub.Fishes {
		if distance(fish.X, fish.Y, req.X, req.Y) <= radius {
			fish.X += (req.X - fish.X) * 0.5
			fish.Y += (req.Y - fish.Y) * 0.5
		}
	}

	msg := &ws.MagnetEffect{
		UserId:  req.UserId,
		CenterX: req.X,
		CenterY: req.Y,
		Radius:  radius,
	}
	data, _ := proto.Marshal(msg)
	wrapped := &ws.WSMessage{
		Event: "magnet_effect",
		Data:  data,
	}
	encoded, _ := proto.Marshal(wrapped)
	c.RoomHub.Broadcast <- encoded

	broadcastSkillUse("magnet", req.UserId, c)
}

func handleEMPBlastSkill(c *Client, req *ws.UseSkillRequest, radius float32, duration time.Duration) {
	c.RoomHub.mu.Lock()
	for _, fish := range c.RoomHub.Fishes {
		if distance(fish.X, fish.Y, req.X, req.Y) <= radius {
			fish.SpeedX = 0
			fish.SpeedY = 0
		}
	}
	c.RoomHub.mu.Unlock()

	go func() {
		time.Sleep(duration)
		c.RoomHub.mu.Lock()
		defer c.RoomHub.mu.Unlock()
		for _, fish := range c.RoomHub.Fishes {
			if distance(fish.X, fish.Y, req.X, req.Y) <= radius {
				// 复活速度（可根据你逻辑调整）
				fish.SpeedX = rand.Float32()*2 - 1
				fish.SpeedY = rand.Float32()*2 - 1
			}
		}
	}()

	msg := &ws.EMPBlastEffect{
		UserId:   req.UserId,
		CenterX:  req.X,
		CenterY:  req.Y,
		Radius:   radius,
		Duration: int32(duration.Seconds()),
	}
	broadcastToAll("emp_blast_effect", msg, c)

	broadcastSkillUse("emp_blast", req.UserId, c)
}

func handleInvisibleSkill(c *Client, req *ws.UseSkillRequest, duration time.Duration) {
	msg := &ws.Invisibility{
		UserId:   req.UserId,
		Duration: int32(duration.Seconds()),
	}
	broadcastToAll("invisibility", msg, c)

	broadcastSkillUse("invisible", req.UserId, c)
}

func handleGrowFishSkill(c *Client, req *ws.UseSkillRequest, scaleFactor float32) {
	c.RoomHub.mu.Lock()
	for _, fish := range c.RoomHub.Fishes {
		if distance(fish.X, fish.Y, req.X, req.Y) <= 100 {
			fish.Size *= scaleFactor
		}
	}
	c.RoomHub.mu.Unlock()

	msg := &ws.FishGrowEffect{
		UserId:  req.UserId,
		CenterX: req.X,
		CenterY: req.Y,
		Radius:  100,
		Scale:   scaleFactor,
	}
	broadcastToAll("fish_grow_effect", msg, c)

	broadcastSkillUse("grow_fish", req.UserId, c)
}
