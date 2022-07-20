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

func TestReverse(t *testing.T) {
	assert.Equal(t, "321ateb", qutil.ReverseString("beta123"))
}

func TestCommonPrefix(t *testing.T) {
	assert.Equal(t, "", qutil.CommonPrefix([]string{"alpha", "beta"}))
	assert.Equal(t, "foo-", qutil.CommonPrefix([]string{"foo-alpha", "foo-beta"}))
	assert.Equal(t, "dc•", qutil.CommonPrefix([]string{"dc•alpha", "dc•beta"}))
}

func TestCommonSuffix(t *testing.T) {
	assert.Equal(t, "", qutil.CommonSuffix([]string{"foo", "bar"}))
	assert.Equal(t, "a-foo", qutil.CommonSuffix([]string{"alpha-foo", "beta-foo"}))
	assert.Equal(t, "a•dc", qutil.CommonSuffix([]string{"alpha•dc", "beta•dc"}))
}

func TestStripQuakeFixes(t *testing.T) {
	t.Run("prefix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"a.alpha", "a.beta"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"dc•alpha", "dc•beta"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"••••alpha", "••••beta"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"••alpha", "••beta"}))
		assert.Equal(t, []string{"alpha", "alphabet"}, qutil.StripQuakeFixes([]string{"dc•alpha", "dc•alphabet"}))
		assert.Equal(t, []string{".alpha", ".beta"}, qutil.StripQuakeFixes([]string{".alpha", ".beta"}))
	})

	t.Run("suffix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha.a", "beta.a"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha•dc", "beta•dc"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha••••", "beta••••"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha••", "beta••"}))
		assert.Equal(t, []string{"alpha.", "beta."}, qutil.StripQuakeFixes([]string{"alpha.", "beta."}))
	})

	t.Run("prefix and suffix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"a.alpha.b", "a.beta.b"}))
		assert.Equal(t, []string{".alpha.", ".beta."}, qutil.StripQuakeFixes([]string{".alpha.", ".beta."}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"||alpha--", "||beta--"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"))alpha((", "))beta(("}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"__alpha..", "__beta.."}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"•••alpha•••", "•••beta•••"}))
	})

	t.Run("short values", func(t *testing.T) {
		assert.Equal(t, []string{".a", ".b"}, qutil.StripQuakeFixes([]string{".a", ".b"}))
	})

	t.Run("single value", func(t *testing.T) {
		assert.Equal(t, []string{"...alpha"}, qutil.StripQuakeFixes([]string{"...alpha"}))
	})
}

func BenchmarkStripQuakeFixes(b *testing.B) {
	names := []string{"•••fox•••", "•••alpha•••", "•••beta•••", "•••delta•••", "•••gamma•••", "•••epsilon•••"}

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("few values", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qutil.StripQuakeFixes(names[0:2])
		}
	})

	b.Run("many values", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qutil.StripQuakeFixes(names)
		}
	})
}

func TestMarshalNoEscapeHtml(t *testing.T) {
	jsonValue, err := qutil.MarshalNoEscapeHtml("<foo&bar>")
	assert.Equal(t, "\"<foo&bar>\"\n", string(jsonValue))
	assert.Nil(t, err)
}

func BenchmarkMarshalNoEscapeHtml(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qutil.MarshalNoEscapeHtml("<foo&bar>")
	}
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

func TestHostnameToIp(t *testing.T) {
	assert.Equal(t, "91.102.91.59", qutil.HostnameToIp("qw.foppa.dk"))
	assert.Equal(t, "91.102.91.59", qutil.HostnameToIp("91.102.91.59"))
}

func BenchmarkHostnameToIp(b *testing.B) {
	b.Run("ip", func(b *testing.B) {
		qutil.HostnameToIp("91.102.91.59")
	})

	b.Run("hostname", func(b *testing.B) {
		qutil.HostnameToIp("qw.foppa.dk")
	})
}

func TestPluralize(t *testing.T) {
	assert.Equal(t, "players", qutil.Pluralize("player", 0))
	assert.Equal(t, "player", qutil.Pluralize("player", 1))
	assert.Equal(t, "players", qutil.Pluralize("player", 2))
}

func TestWildcardMatchStringSlice(t *testing.T) {
	const wildcard = "@"
	assert.False(t, qutil.WildcardMatchStringSlice(nil, "foo", wildcard))
	assert.False(t, qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "gamma"}, "foo", wildcard))
	assert.True(t, qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "foo", "gamma"}, "foo", wildcard))
	assert.True(t, qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "foo", "gamma"}, "FOO", wildcard))
	assert.True(t, qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "FOO", "gamma"}, "foo", wildcard))
	assert.True(t, qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "FOO", "gamma"}, "@amma", wildcard))
	assert.True(t, qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "FOO", "gamma"}, "gamm@", wildcard))
	assert.True(t, qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "FOO", "gamma"}, "@amm@", wildcard))
}

func BenchmarkWildcardMatchStringSlice(b *testing.B) {
	b.ReportAllocs()
	const wildcard = "@"

	for i := 0; i < b.N; i++ {
		qutil.WildcardMatchStringSlice([]string{"alpha", "beta", "foo", "gamma"}, "foo", wildcard)
	}
}

func TestWildcardMatchString(t *testing.T) {
	const wildcard = "@"
	assert.False(t, qutil.WildcardMatchString("", "beta", wildcard))
	assert.False(t, qutil.WildcardMatchString("alpha", "", wildcard))
	assert.False(t, qutil.WildcardMatchString("alpha", "beta", wildcard))

	assert.True(t, qutil.WildcardMatchString("alpha", "alpha", wildcard))
	assert.True(t, qutil.WildcardMatchString("ALPHA", "alpha", wildcard))
	assert.True(t, qutil.WildcardMatchString("alpha", "ALPHA", wildcard))
	assert.True(t, qutil.WildcardMatchString("ALPHA", "ALPHA", wildcard))

	// prefix wildcard
	assert.True(t, qutil.WildcardMatchString("alphabetic", "@betic", wildcard))
	assert.True(t, qutil.WildcardMatchString("betic", "@betic", wildcard))

	// suffix wilcard
	assert.True(t, qutil.WildcardMatchString("alphabetic", "alpha@", wildcard))
	assert.True(t, qutil.WildcardMatchString("alpha", "alpha@", wildcard))

	// suffix and prefix wildcard
	assert.True(t, qutil.WildcardMatchString("alphabetic", "@lphabeti@", wildcard))
	assert.True(t, qutil.WildcardMatchString("alphabetic", "@alphabetic@", wildcard))
}
