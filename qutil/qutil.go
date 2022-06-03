package qutil

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/goccy/go-json"
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
	longestPrefix := strings.Builder{}

	if len(strs) > 0 {
		sort.Strings(strs)
		firstRunes := []rune(strs[0])
		lastRunes := []rune(strs[len(strs)-1])

		for i := 0; i < len(firstRunes); i++ {
			if lastRunes[i] == firstRunes[i] {
				longestPrefix.WriteRune(lastRunes[i])
			} else {
				break
			}
		}
	}

	return longestPrefix.String()
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

	// minimum fix to strip
	minFixLength := 2

	// skip if any value is equal to or shorter than min fix length
	for _, value := range strs {
		if len(value) <= minFixLength {
			return strs
		}
	}

	const delimiterChars = ".•_-|[]{}()"

	// prefix
	prefix := CommonPrefix(strs)

	if strings.ContainsAny(prefix, delimiterChars) {
		lastDelimiterIndex := strings.LastIndexAny(prefix, delimiterChars)
		prefixRuneCount := len([]rune(prefix))

		// utf8 runes have length > 1
		if lastDelimiterIndex > prefixRuneCount {
			lastDelimiterIndex = prefixRuneCount - 1
		}

		prefixLength := lastDelimiterIndex + 1

		if prefixLength >= minFixLength {
			for index := range strs {
				runes := []rune(strs[index])
				strs[index] = string(runes[prefixLength:])
			}
		}
	}

	// suffix
	suffix := CommonSuffix(strs)

	if strings.ContainsAny(suffix, delimiterChars) {
		firstDelimiterIndex := strings.IndexAny(suffix, delimiterChars)
		suffixLength := len([]rune(suffix)) - firstDelimiterIndex

		if suffixLength >= minFixLength {
			for index := range strs {
				runes := []rune(strs[index])
				newLastIndex := len(runes) - suffixLength
				strs[index] = string(runes[0:newLastIndex])
			}
		}
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

func TrimSymbols(value string) string {
	result := strings.Builder{}

	lastWasChar := false
	lastWasDelimiter := false

	for _, r := range []rune(strings.Trim(value, " -")) {
		if (' ' == r || '-' == r) && !lastWasDelimiter {
			result.WriteRune(r)
			lastWasChar = false
			lastWasDelimiter = true
		} else if unicode.IsOneOf([]*unicode.RangeTable{unicode.Letter, unicode.Number}, r) {
			result.WriteRune(r)
			lastWasChar = true
			lastWasDelimiter = false
		} else if lastWasChar {
			result.WriteRune(' ')
			lastWasChar = false
			lastWasDelimiter = false
		}
	}

	return strings.TrimRight(result.String(), " ")
}
