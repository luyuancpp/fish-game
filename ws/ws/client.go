package ws

import (
	"fish-game/apps/room/room"
	"fish-game/apps/user/user"
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	UserID      string
	RoomID      string
	Conn        *websocket.Conn
	Send        chan []byte
	RoomHub     *RoomHub
	RoomClient  room.RoomClient
	UserClient  user.UserClient // ✅ 加这一
	Cooldowns   map[string]time.Time
	SkillLevels map[string]int

	// ✅ 新增字段：背包（道具库存）
	Inventory map[string]int
}

func (c *Client) GetUserID() string {
	return c.UserID
}
