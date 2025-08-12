package ws

import (
	"fish-game/apps/room/room"
	"fish-game/apps/user/user"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID     string
	RoomID     string
	Conn       *websocket.Conn
	Send       chan []byte
	RoomHub    *RoomHub
	RoomClient room.RoomClient
	UserClient user.UserClient // ✅ 加这一
}

func (c *Client) GetUserID() string {
	return c.UserID
}
