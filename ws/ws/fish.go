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
			data := map[string]interface{}{
				"event": "fish_generate",
				"data": map[string]interface{}{
					"id": fishID,
					"x":  100,
					"y":  50,
				},
			}
			msg, _ := json.Marshal(data)
			hub.Broadcast <- msg
		}
	}()
}
