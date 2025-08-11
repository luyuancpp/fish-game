package ws

import (
	"fish-game/apps/room/room"
	"fish-game/apps/user/user"
	"fish-game/pkg/jwt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var roomHubs = make(map[string]*RoomHub)

func getOrCreateRoom(roomID string) *RoomHub {
	// 判断房间是否已被其他实例绑定
	targetInstance, err := GetRoomWS(roomID)
	if err == nil && targetInstance != localWSID {
		log.Printf("⚠️ 房间 %s 已绑定到实例 %s，本实例是 %s，拒绝处理", roomID, targetInstance, localWSID)
		return nil // 标志性拒绝
	}

	// 如果已存在本地 hub
	if hub, ok := roomHubs[roomID]; ok {
		return hub
	}

	// 绑定房间到当前实例
	err = BindRoomToWS(roomID)
	if err != nil {
		log.Printf("❌ 房间绑定失败: %v", err)
	}

	// 创建 RoomHub
	hub := NewRoomHub(roomID)
	roomHubs[roomID] = hub
	go hub.Run()
	StartFishGenerator(hub)
	return hub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	RoomClient room.RoomClient
	UserClient user.UserClient // ✅ 加这一行
}

func NewWSHandler(roomClient room.RoomClient, userClient user.UserClient) *WSHandler {
	return &WSHandler{
		RoomClient: roomClient,
		UserClient: userClient,
	}
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// 1. 验证 token
	token := r.URL.Query().Get("token")
	uid, ok := jwtx.VerifyToken(token, jwtx.DefaultSecret)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. 调用 Room RPC 获取房间号
	reply, err := h.RoomClient.JoinRoom(r.Context(), &room.JoinRoomRequest{
		Uid: strconv.FormatInt(uid, 10),
	})
	if err != nil {
		log.Println("❌ JoinRoom RPC error:", err)
		http.Error(w, "JoinRoom failed", http.StatusInternalServerError)
		return
	}
	roomId := reply.RoomId
	log.Printf("✅ 玩家 %d 加入房间 %s", uid, roomId)

	// 3. 升级 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	// 4. 加入 RoomHub（自动创建或复用）
	hub := getOrCreateRoom(roomId)
	if hub == nil {
		log.Printf("🚫 拒绝加入房间 %s，因本实例未绑定该房间", roomId)
		http.Error(w, "Room is handled by another server", http.StatusServiceUnavailable)
		return
	}

	client := &Client{
		UserID:  strconv.FormatInt(uid, 10),
		RoomID:  roomId,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		RoomHub: hub,
	}

	// 注册客户端
	hub.Register <- client

	go client.writePump()
	go client.readPump()
}
