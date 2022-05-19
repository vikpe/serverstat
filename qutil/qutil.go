package qutil

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
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

func StripQuakeFixes(strs []string) []string {
	// skip if few values
	if len(strs) < 2 {
		return strs
	}

	// minimium fix to strip
	minFixLength := 2

	// skip if any value is equal to or shorter than min fix length
	for _, value := range strs {
		if len(value) <= minFixLength {
			return strs
		}
	}

	delimiterChars := ".•_-|[]{}()"

	// prefix
	prefix := CommonPrefix(strs)

	if strings.ContainsAny(prefix, delimiterChars) {
		quakePrefix := prefix[0 : strings.LastIndexAny(prefix, delimiterChars)+1]

		if len(quakePrefix) >= minFixLength {
			strs = StripPrefix(strs, quakePrefix)
		}
	}

	// suffix
	suffix := CommonSuffix(strs)

	if strings.ContainsAny(suffix, delimiterChars) {
		quakeSuffix := suffix[strings.IndexAny(suffix, delimiterChars):]

		if len(quakeSuffix) >= minFixLength {
			strs = StripSuffix(strs, quakeSuffix)
		}
	}

	return strs
}

func StripSuffix(strs []string, value string) []string {
	for index, val := range strs {
		strs[index] = strings.TrimSuffix(val, value)
	}

	return strs
}

func StripPrefix(strs []string, value string) []string {
	for index, val := range strs {
		strs[index] = strings.TrimPrefix(val, value)
	}

	return strs
}

func MarshalNoEscapeHtml(value any) ([]byte, error) {
	var dst bytes.Buffer
	enc := json.NewEncoder(&dst)
	enc.SetEscapeHTML(false)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}
