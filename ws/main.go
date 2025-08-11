package main

import (
	"fish-game/apps/user/user"
	"fish-game/ws/ws"
	"log"
	"net/http"
	"os"

	"fish-game/apps/room/room"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	// 1. 启动 WebSocket Hub
	// （RoomHub 是按需自动创建）
	os.Setenv("WS_INSTANCE_ID", "ws-1") // 你也可以通过启动脚本注入

	ws.StartGlobalBroadcastListener()

	// 2. 创建 roomRpc 客户端
	roomClient := room.NewRoomClient(
		zrpc.MustNewClient(zrpc.RpcClientConf{
			Target: "127.0.0.1:8083", // 修改为你的 room-rpc 地址
		}).Conn(),
	)

	userClient := user.NewUserClient(
		zrpc.MustNewClient(zrpc.RpcClientConf{
			Target: "127.0.0.1:8084", // 改成你的 user-rpc 地址
		}).Conn(),
	)

	// 3. 创建 WSHandler，注入 RoomClient
	handler := ws.NewWSHandler(roomClient, userClient)

	http.HandleFunc("/ws", handler.ServeWS)

	log.Println("WebSocket started on :8082")
	http.ListenAndServe(":8082", nil)
}
