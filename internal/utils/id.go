package utils

import (
	"strconv"
)

func IsValidInt64(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}
