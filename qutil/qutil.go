package qutil

import "strconv"

func StringToInt(value string) int {
	valueAsInt, _ := strconv.Atoi(value)
	return valueAsInt
}
