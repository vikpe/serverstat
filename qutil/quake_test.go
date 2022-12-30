package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qutil"
)

func TestStripQuakeFixes(t *testing.T) {
	t.Run("prefix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"a.alpha", "a.beta"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"dc•alpha", "dc•beta"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"••••alpha", "••••beta"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"••alpha", "••beta"}))
		assert.Equal(t, []string{"alpha", "alphabet"}, qutil.StripQuakeFixes([]string{"dc•alpha", "dc•alphabet"}))
		assert.Equal(t, []string{".alpha", ".beta"}, qutil.StripQuakeFixes([]string{".alpha", ".beta"}))
		assert.Equal(t, []string{"---a", "---"}, qutil.StripQuakeFixes([]string{"---a", "---"}))
	})

	t.Run("suffix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha.a", "beta.a"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha•dc", "beta•dc"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha••••", "beta••••"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"alpha••", "beta••"}))
		assert.Equal(t, []string{"alpha.", "beta."}, qutil.StripQuakeFixes([]string{"alpha.", "beta."}))
		assert.Equal(t, []string{"a---", "---"}, qutil.StripQuakeFixes([]string{"a---", "---"}))
	})

	t.Run("prefix and suffix", func(t *testing.T) {
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"a.alpha.b", "a.beta.b"}))
		assert.Equal(t, []string{".alpha.", ".beta."}, qutil.StripQuakeFixes([]string{".alpha.", ".beta."}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"||alpha--", "||beta--"}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"))alpha((", "))beta(("}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"__alpha..", "__beta.."}))
		assert.Equal(t, []string{"alpha", "beta"}, qutil.StripQuakeFixes([]string{"•••alpha•••", "•••beta•••"}))
		assert.Equal(t, []string{"•••a•••", "•••"}, qutil.StripQuakeFixes([]string{"•••a•••", "•••"}))
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
