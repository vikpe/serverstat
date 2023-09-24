package qutil

import (
	"strings"

	"github.com/jpillora/longestcommon"
)

const delimiterChars = `."_-|[]{}():`
const minFixLength = 2

func StripQuakeFixes(strs []string) []string {
	// skip if few values
	if len(strs) < 2 {
		return strs
	}

	// skip if any value is equal to or shorter than min fix length
	for _, value := range strs {
		if len(value) <= minFixLength {
			return strs
		}
	}

	// replace utf8
	for index := range strs {
		strs[index] = strings.ReplaceAll(strs[index], "•", `"`)
	}

	// prefix
	prefix := getPrefix(strs)

	if len(prefix) >= minFixLength {
		for index := range strs {
			strs[index] = strings.TrimPrefix(strs[index], prefix)
		}
	}

	// suffix
	suffix := getSuffix(strs)

	if len(suffix) >= minFixLength {
		for index := range strs {
			strs[index] = strings.TrimSuffix(strs[index], suffix)
		}
	}

	// restore utf8
	for index := range strs {
		strs[index] = strings.ReplaceAll(strs[index], `"`, "•")
	}

	// trim spaces
	for index := range strs {
		strs[index] = strings.TrimSpace(strs[index])
	}

	return strs
}

func getPrefix(strs []string) string {
	fullPrefix := longestcommon.Prefix(strs)

	if !strings.ContainsAny(fullPrefix, delimiterChars) {
		return ""
	}

	lastDelimiterIndex := strings.LastIndexAny(fullPrefix, delimiterChars)
	prefix := fullPrefix[0 : lastDelimiterIndex+1]
	prefixLength := len(prefix)

	for index := range strs {
		if len(strs[index]) <= prefixLength {
			return ""
		}
	}

	return prefix
}

func getSuffix(strs []string) string {
	fullSuffix := longestcommon.Suffix(strs)

	if !strings.ContainsAny(fullSuffix, delimiterChars) {
		return ""
	}

	firstDelimiterIndex := strings.IndexAny(fullSuffix, delimiterChars)
	suffix := fullSuffix[firstDelimiterIndex:]
	suffixLength := len(suffix)

	for index := range strs {
		if len(strs[index]) <= suffixLength {
			return ""
		}
	}

	return suffix
}
