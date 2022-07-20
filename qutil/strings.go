package qutil

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
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

func WildcardMatchStringSlice(haystack []string, needle string, wildcard string) bool {
	if 0 == len(haystack) {
		return false
	}

	for _, value := range haystack {
		if WildcardMatchString(value, needle, wildcard) {
			return true
		}
	}

	return false
}

func WildcardMatchString(haystack, needle, wildcard string) bool {
	if strings.Contains(needle, wildcard) {
		hasWildcardPrefix := strings.HasPrefix(needle, wildcard)
		hasWildCardSuffix := strings.HasSuffix(needle, wildcard)
		haystackLower := strings.ToLower(haystack)
		needleLower := strings.ToLower(needle)

		if hasWildcardPrefix && hasWildCardSuffix {
			return strings.Contains(haystackLower, needleLower[1:len(needleLower)-1])
		} else if hasWildcardPrefix {
			return strings.HasSuffix(haystackLower, needleLower[1:])
		} else { // hasWildcardPrefix
			return strings.HasPrefix(haystackLower, needleLower[:len(needleLower)-1])
		}
	}

	if len(needle) != len(haystack) {
		return false
	}

	return strings.EqualFold(needle, haystack)
}
