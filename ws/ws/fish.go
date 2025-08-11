package ws

import (
	"encoding/json"
	"time"
)

func StartFishGenerator(hub *RoomHub) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		fishID := 0
		for range ticker.C {
			fishID++

			fish := &pb.FishGenerate{
				Id: int32(fishID),
				X:  rand.Float32() * 500,
				Y:  rand.Float32() * 300,
			}

			data, _ := proto.Marshal(fish)

			msg := &pb.WSMessage{
				Event: "fish_generate",
				Data:  data,
			}

			encoded, _ := proto.Marshal(msg)
			hub.Broadcast <- encoded
		}
	}()
}
