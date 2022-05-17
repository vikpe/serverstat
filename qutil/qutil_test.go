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
}

func TestCommonSuffix(t *testing.T) {
	assert.Equal(t, "", qutil.CommonSuffix([]string{"foo", "bar"}))
	assert.Equal(t, "a-foo", qutil.CommonSuffix([]string{"alpha-foo", "beta-foo"}))
}

func TestStripSuffix(t *testing.T) {
	assert.Equal(t, []string{"alpha", "beta"}, qutil.StripSuffix([]string{"alpha-foo", "beta-foo"}, "-foo"))
}

func TestStripPrefix(t *testing.T) {
	assert.Equal(t, []string{"alpha", "beta"}, qutil.StripPrefix([]string{"foo-alpha", "foo-beta"}, "foo-"))
}

func TestStripQuakeFixes(t *testing.T) {
	t.Run("prefix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"a.alpha", "a.beta"}))
		assert.Equal(t, []string{".alpha", ".beta"}, qutil.StripQuakeFixes([]string{".alpha", ".beta"}))
	})

	t.Run("suffix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha.a", "beta.a"}))
		assert.Equal(t, []string{"alpha.", "beta."}, qutil.StripQuakeFixes([]string{"alpha.", "beta."}))
	})

	t.Run("prefix and suffix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"a.alpha.b", "a.beta.b"}))
		assert.Equal(t, []string{".alpha.", ".beta."}, qutil.StripQuakeFixes([]string{".alpha.", ".beta."}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"||alpha--", "||beta--"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"))alpha((", "))beta(("}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"__alpha..", "__beta.."}))
	})

	t.Run("short values", func(t *testing.T) {
		assert.Equal(t, []string{".a", ".b"}, qutil.StripQuakeFixes([]string{".a", ".b"}))
	})

	t.Run("single value", func(t *testing.T) {
		assert.Equal(t, []string{"...alpha"}, qutil.StripQuakeFixes([]string{"...alpha"}))
	})
}
