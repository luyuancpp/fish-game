package ws

import (
	"fish-game/apps/room/room"
	"fish-game/pkg/jwt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var roomHubs = make(map[string]*RoomHub)

func getOrCreateRoom(roomId string) *RoomHub {
	if hub, ok := roomHubs[roomId]; ok {
		return hub
	}
	hub := NewRoomHub(roomId)
	roomHubs[roomId] = hub
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
}

func NewWSHandler(client room.RoomClient) *WSHandler {
	return &WSHandler{
		RoomClient: client,
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
