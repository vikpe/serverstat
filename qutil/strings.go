package qutil

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/IGLOU-EU/go-wildcard"
)

func StringToInt(value string) int {
	valueAsInt, _ := strconv.Atoi(value)
	return valueAsInt
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

func Pluralize(word string, count int) string {
	if 1 == count {
		return word
	}

	return fmt.Sprintf("%ss", word)
}

func WildcardMatchStringSlice(haystack []string, pattern string) bool {
	if 0 == len(haystack) {
		return false
	}

	for _, value := range haystack {
		if wildcard.MatchSimple(strings.ToLower(pattern), strings.ToLower(value)) {
			return true
		}
	}

	return false
}
