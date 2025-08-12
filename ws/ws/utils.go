package ws

import (
	"strconv"
)

func AtoiSafe(s string) int32 {
	i, _ := strconv.Atoi(s)
	return int32(i)
}
