package qstring_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestQuakeString_ToPlainString(t *testing.T) {
	testCases := map[string]string{
		"HCBtYXplcg==": "• mazer",
		"EXNyEA==":     "]sr[",
		"W1NlcnZlTWVd": "[ServeMe]",
		"bHF3Yw==":     "lqwc",
		"4uHz8w==":     "bass",
	}

	for encodedString, expect := range testCases {
		strBytes, _ := base64.StdEncoding.DecodeString(encodedString)
		actual := qstring.New(string(strBytes)).ToPlainString()
		assert.Equal(t, expect, actual)
	}
}

func TestQuakeString_ToColorCodes(t *testing.T) {
	testCases := map[string]string{
		"HCBtYXplcg==": "wwwwwww", // • mazer
		"EXNyEA==":     "gwwg",    // ]sr[
		"4uHz8w==":     "bbbb",    // bass
	}

	for encodedString, expect := range testCases {
		strBytes, _ := base64.StdEncoding.DecodeString(encodedString)
		actual := qstring.New(string(strBytes)).ToColorCodes()
		assert.Equal(t, expect, actual)
	}
}

func TestQuakeString_MarshalJSON(t *testing.T) {
	jsonValue, _ := qstring.New("<XantoM>").MarshalJSON()
	assert.Equal(t, "\"<XantoM>\"\n", string(jsonValue))
}
