package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qutil"
)

func TestStringToInt(t *testing.T) {
	assert.Equal(t, 1, qutil.StringToInt("1"))
	assert.Equal(t, 55, qutil.StringToInt("55"))
	assert.Equal(t, -55, qutil.StringToInt("-55"))
	assert.Equal(t, 0, qutil.StringToInt("0"))
	assert.Equal(t, 0, qutil.StringToInt(""))
}

func TestTrimSymbols(t *testing.T) {
	var testCases = map[string]string{
		"--great--":                      "great",
		"--great-stuff--":                "great-stuff",
		"***":                            "",
		"..ab":                           "ab",
		"..abc._[]{}()*\\/def*--2on2* -": "abc def -2on2",
		"..abc":                          "abc",
		".a.b. ":                         "a b",
		".a-b. ":                         "a-b",
		"[2on2]":                         "2on2",
		"*race...":                       "race",
	}

	for value, expect := range testCases {
		assert.Equal(t, expect, qutil.TrimSymbols(value))
	}
}

func BenchmarkTrimSymbols(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("short value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qutil.TrimSymbols("-.a.b. ")
		}
	})

	b.Run("long value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qutil.TrimSymbols("..acb._[]{}()*\\/def*2on2* -")
		}
	})
}

func TestPluralize(t *testing.T) {
	assert.Equal(t, "players", qutil.Pluralize("player", 0))
	assert.Equal(t, "player", qutil.Pluralize("player", 1))
	assert.Equal(t, "players", qutil.Pluralize("player", 2))
}

func TestStripControlCharacters(t *testing.T) {
	input := string("\tfoo\fbar\n\a\b \vbaz")
	expected := string("foobar baz")
	actual := qutil.StripControlCharacters(input)
	assert.Equal(t, expected, actual)
}

func TestUnicodeToAscii(t *testing.T) {
	assert.Equal(t, "XantoM", qutil.UnicodeToAscii("XantoM"))
	assert.Equal(t, "nigve", qutil.UnicodeToAscii("\u00EE\u00E9\u00E7\u00F6\u00E5"))
}
