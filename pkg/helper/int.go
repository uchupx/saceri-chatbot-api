package helper

import (
	"strconv"
)

func StringToUint(s string) uint {
	u, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(u)
}

func DefaultInt(v *int, def int) int {
	if v == nil {
		return def
	}

	return *v
}
