package ws

import "math"

func hit(fish *Fish, bulletX, bulletY float32) bool {
	dx := fish.X - bulletX
	dy := fish.Y - bulletY
	distance := float32(math.Hypot(float64(dx), float64(dy)))
	return distance <= fish.HitRadius
}
