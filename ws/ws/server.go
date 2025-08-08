package ws

import (
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

func ServeWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	roomId := r.URL.Query().Get("roomId")

	userId, ok := jwtx.VerifyToken(token, jwtx.DefaultSecret)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	hub := getOrCreateRoom(roomId)
	client := &Client{
		UserID:  strconv.FormatInt(userId, 10),
		RoomID:  roomId,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		RoomHub: hub,
	}
	hub.Register <- client

	go client.writePump()
	go client.readPump()
}
