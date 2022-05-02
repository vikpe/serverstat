package qstring_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-serverstat/quaketext/qstring"
)

func TestToPlainString(t *testing.T) {
	testCases := map[string]string{
		"HCBtYXplcg==": "• mazer",
		"EXNyEA==":     "]sr[",
		"W1NlcnZlTWVd": "[ServeMe]",
		"bHF3Yw==":     "lqwc",
		"4uHz8w==":     "bass",
	}

	for encodedString, expect := range testCases {
		strBytes, _ := base64.StdEncoding.DecodeString(encodedString)
		actual := qstring.ToPlainString(string(strBytes))
		assert.Equal(t, expect, actual)
	}
}
