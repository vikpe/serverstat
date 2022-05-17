package qutil

import (
	"sort"
	"strconv"
)

func StringToInt(value string) int {
	valueAsInt, _ := strconv.Atoi(value)
	return valueAsInt
}

func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func CommonPrefix(strs []string) string {
	var longestPrefix = ""

	if len(strs) > 0 {
		sort.Strings(strs)
		first := strs[0]
		last := strs[len(strs)-1]

		for i := 0; i < len(first); i++ {
			if string(last[i]) == string(first[i]) {
				longestPrefix += string(last[i])
			} else {
				break
			}
		}
	}
	return longestPrefix
}

func CommonSuffix(strs []string) string {
	reversedStrings := make([]string, 0)

	for _, s := range strs {
		reversedStrings = append(reversedStrings, ReverseString(s))
	}

	return ReverseString(CommonPrefix(reversedStrings))
}
