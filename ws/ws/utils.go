package ws

import (
	"strconv"
)

func atoiSafe(s string) int32 {
	i, _ := strconv.Atoi(s)
	return int32(i)
}
