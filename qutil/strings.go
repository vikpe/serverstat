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

func StripControlCharacters(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r != 127 {
			return r
		}
		return -1
	}, str)
}

func UnicodeToAscii(str string) string {
	result := strings.Builder{}
	const AsciiMax = 128

	for _, r := range str {
		if r >= AsciiMax {
			result.WriteRune(r - AsciiMax)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}
