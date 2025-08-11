package ws

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr: "127.0.0.1:6379", // 改成你的 Redis 地址
})

var localWSID = os.Getenv("WS_INSTANCE_ID") // 实例启动时指定

// 把房间绑定到当前实例
func BindRoomToWS(roomID string) error {
	key := fmt.Sprintf("room_ws:%s", roomID)
	return redisClient.Set(ctx, key, localWSID, 0).Err()
}

// 获取房间绑定的实例
func GetRoomWS(roomID string) (string, error) {
	key := fmt.Sprintf("room_ws:%s", roomID)
	return redisClient.Get(ctx, key).Result()
}

// 监听来自其他 WS 实例的广播
func StartGlobalBroadcastListener() {
	pubsub := redisClient.Subscribe(ctx, "global_broadcast")
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			for _, hub := range roomHubs {
				hub.Broadcast <- []byte(msg.Payload)
			}
		}
	}()
}

// 发送跨实例广播
func GlobalBroadcast(data []byte) error {
	return redisClient.Publish(ctx, "global_broadcast", data).Err()
}
