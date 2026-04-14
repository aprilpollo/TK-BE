package utils

import (
	"strconv"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidInt64(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}