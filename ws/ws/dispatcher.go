package ws

import (
	ws "fish-game/ws/proto"
)

var handlers = make(map[string]func(c *Client, msg *ws.WSMessage))

func init() {
	handlers["shoot"] = HandleShoot
	handlers["use_skill"] = HandleUseSkill
}
