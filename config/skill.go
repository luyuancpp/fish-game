package config

import "time"

type SkillConfig struct {
	Cooldown       time.Duration
	RequiredItem   string
	BaseDuration   time.Duration // 适用于如 freeze 等技能
	BaseMultiplier float64       // 适用于 grow_fish 等技能
	// 你可以扩展更多字段，比如作用范围、技能图标、动画ID等
}

var SkillConfigs = map[string]SkillConfig{
	"freeze": {
		Cooldown:     10 * time.Second,
		RequiredItem: "item_freeze",
		BaseDuration: 5 * time.Second,
	},
	"missile": {
		Cooldown:     12 * time.Second,
		RequiredItem: "item_missile",
	},
	"grow_fish": {
		Cooldown:       20 * time.Second,
		RequiredItem:   "item_grow",
		BaseMultiplier: 2.0,
	},
	// 更多技能...
}
