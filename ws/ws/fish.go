package ws

import (
	ws "fish-game/ws/proto"
	"google.golang.org/protobuf/proto"
	"math"
	"math/rand"
	"time"
)

// 鱼结构体定义
type Fish struct {
	Id        int32
	X, Y      float32
	SpeedX    float32
	SpeedY    float32
	HitRadius float32
}

// 生成一条随机初始位置和速度的鱼
func NewFish(id int32) *Fish {
	angle := rand.Float64() * 2 * math.Pi
	speed := float32(20 + rand.Intn(30)) // 速度范围: 20~50 单位/秒

	return &Fish{
		Id:        id,
		X:         rand.Float32() * 500,
		Y:         rand.Float32() * 300,
		SpeedX:    float32(math.Cos(angle)) * speed,
		SpeedY:    float32(math.Sin(angle)) * speed,
		HitRadius: 30,
	}
}

// 鱼移动逻辑，dt 为移动间隔（单位秒）
func (f *Fish) Move(dt float32) {
	f.X += f.SpeedX * dt
	f.Y += f.SpeedY * dt

	// 边界检测与反弹
	if f.X < 0 || f.X > 500 {
		f.SpeedX = -f.SpeedX
	}
	if f.Y < 0 || f.Y > 300 {
		f.SpeedY = -f.SpeedY
	}

	// 再次裁剪坐标防止越界
	if f.X < 0 {
		f.X = 0
	}
	if f.X > 500 {
		f.X = 500
	}
	if f.Y < 0 {
		f.Y = 0
	}
	if f.Y > 300 {
		f.Y = 300
	}
}

// 启动鱼生成 + 移动 + 广播的循环
func StartFishGenerator(hub *RoomHub) {
	go func() {
		tickerGenerate := time.NewTicker(3 * time.Second)    // 每3秒生成一条新鱼
		tickerMove := time.NewTicker(100 * time.Millisecond) // 每0.1秒移动鱼并同步
		defer tickerGenerate.Stop()
		defer tickerMove.Stop()

		fishID := int32(0)

		for {
			select {
			case <-tickerGenerate.C:
				// 生成新鱼并广播
				fishID++
				fish := NewFish(fishID)

				hub.mu.Lock()
				hub.Fishes = append(hub.Fishes, fish)
				hub.mu.Unlock()

				msg := &ws.FishGenerate{
					Id: fish.Id,
					X:  fish.X,
					Y:  fish.Y,
				}
				data, _ := proto.Marshal(msg)
				wsMsg := &ws.WSMessage{
					Event: "fish_generate",
					Data:  data,
				}
				encoded, _ := proto.Marshal(wsMsg)
				hub.Broadcast <- encoded

			case <-tickerMove.C:
				// 移动所有鱼，并广播位置信息
				hub.mu.Lock()

				// 👉 新增：判断是否被冻结
				if !hub.IsFrozen() {
					dt := float32(0.1) // 100ms = 0.1秒
					for _, fish := range hub.Fishes {
						fish.Move(dt)
					}
				}

				// ✅ 无论是否冻结，都要广播当前位置（客户端才知道鱼停住了）
				positions := make([]*ws.FishPosition, 0, len(hub.Fishes))
				for _, fish := range hub.Fishes {
					positions = append(positions, &ws.FishPosition{
						Id: fish.Id,
						X:  fish.X,
						Y:  fish.Y,
					})
				}

				posUpdate := &ws.FishPositionUpdate{
					Positions: positions,
				}
				data, _ := proto.Marshal(posUpdate)
				wsMsg := &ws.WSMessage{
					Event: "fish_position_update",
					Data:  data,
				}
				encoded, _ := proto.Marshal(wsMsg)
				hub.Broadcast <- encoded

				hub.mu.Unlock()

			}
		}
	}()
}
