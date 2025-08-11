package ws

import (
	ws "fish-game/ws/proto"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"time"
)

type Fish struct {
	Id        int32
	X, Y      float32
	HitRadius float32
	// 可以后续扩展速度、方向等字段
}

func StartFishGenerator(hub *RoomHub) {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		fishID := 0

		for range ticker.C {
			fishID++

			// 1. 创建服务端鱼对象
			fish := &Fish{
				Id:        fishID,
				X:         rand.Float32() * 500,
				Y:         rand.Float32() * 300,
				HitRadius: 30, // 可调整
			}

			// 2. 存入 RoomHub
			hub.Fishes = append(hub.Fishes, fish)

			// 3. 构造 FishGenerate 消息（用于广播）
			msg := &ws.FishGenerate{
				Id: int32(fish.Id),
				X:  fish.X,
				Y:  fish.Y,
			}
			data, _ := proto.Marshal(msg)

			// 4. 包装成 WSMessage
			wsMsg := &ws.WSMessage{
				Event: "fish_generate",
				Data:  data,
			}
			encoded, _ := proto.Marshal(wsMsg)

			// 5. 广播给所有客户端
			hub.Broadcast <- encoded
		}
	}()
}
