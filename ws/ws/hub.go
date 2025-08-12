package ws

import (
	"sync"
	"time"
)

type RoomHub struct {
	RoomID      string
	Clients     map[*Client]bool
	Broadcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
	Fishes      []*Fish        // 当前房间的所有鱼
	PlayerCoins map[string]int // 玩家金币记录
	mu          sync.Mutex     // ✅ 加上这个字段
	FreezeUntil time.Time      // 冻结到的时间
}

func NewRoomHub(roomId string) *RoomHub {
	return &RoomHub{
		RoomID:     roomId,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *RoomHub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true
		case c := <-h.Unregister:
			delete(h.Clients, c)
			close(c.Send)
		case msg := <-h.Broadcast:
			for c := range h.Clients {
				c.Send <- msg
			}
		}
	}
}

func (h *RoomHub) IsFrozen() bool {
	return time.Now().Before(h.FreezeUntil)
}

func (h *RoomHub) SetFreeze(duration time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.FreezeUntil = time.Now().Add(duration)
}
