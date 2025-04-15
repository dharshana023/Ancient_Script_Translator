package utils

import (
	"strconv"
)

// StringToInt converts a string to an integer
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}