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
