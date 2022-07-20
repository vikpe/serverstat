package qutil

import "strings"

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
