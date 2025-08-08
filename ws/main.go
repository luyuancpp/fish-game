package main

import (
	"fish-game/ws/ws"
	"log"
	"net/http"

	"fish-game/apps/room/room"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	// 1. 启动 WebSocket Hub
	// （RoomHub 是按需自动创建）

	// 2. 创建 roomRpc 客户端
	roomClient := room.NewRoomClient(
		zrpc.MustNewClient(zrpc.RpcClientConf{
			Target: "127.0.0.1:9001", // 修改为你的 room-rpc 地址
		}).Conn(),
	)

	// 3. 创建 WSHandler，注入 RoomClient
	handler := ws.NewWSHandler(roomClient)

	http.HandleFunc("/ws", handler.ServeWS)

	log.Println("WebSocket started on :8082")
	http.ListenAndServe(":8082", nil)
}
