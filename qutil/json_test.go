package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qutil"
)

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
