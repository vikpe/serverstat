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

func BenchmarkToPlainString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	encodedStrings := []string{
		"HCBtYXplcg==", //"• mazer",
		"EXNyEA==",     //    "]sr[",
		"W1NlcnZlTWVd", //"[ServeMe]",
		"bHF3Yw==",     //    "lqwc",
		"4uHz8w==",     //    "bass",
	}

	testStrings := make([]string, 0)

	for _, encodedString := range encodedStrings {
		strBytes, _ := base64.StdEncoding.DecodeString(encodedString)
		testStrings = append(testStrings, string(strBytes))
	}

	for i := 0; i < b.N; i++ {
		for _, str := range testStrings {
			qstring.ToPlainString(str)
		}
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

func BenchmarkToColorCodes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	encodedStrings := []string{
		"HCBtYXplcg==", //"• mazer",
		"EXNyEA==",     //    "]sr[",
		"W1NlcnZlTWVd", //"[ServeMe]",
		"bHF3Yw==",     //    "lqwc",
		"4uHz8w==",     //    "bass",
	}

	testStrings := make([]string, 0)

	for _, encodedString := range encodedStrings {
		strBytes, _ := base64.StdEncoding.DecodeString(encodedString)
		testStrings = append(testStrings, string(strBytes))
	}

	for i := 0; i < b.N; i++ {
		for _, str := range testStrings {
			qstring.ToColorCodes(str)
		}
	}
}

func TestQuakeString_MarshalJSON(t *testing.T) {
	jsonValue, _ := qstring.New("<Xan&toM>").MarshalJSON()
	assert.Equal(t, "\"<Xan&toM>\"\n", string(jsonValue))
}
