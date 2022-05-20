package qtvstream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
)

func TestQtvStream_Url(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		stream := qtvstream.QtvStream{Id: 0, Address: ""}
		assert.Equal(t, "", stream.Url())
	})

	t.Run("not empty", func(t *testing.T) {
		stream := qtvstream.QtvStream{Id: 12, Address: "qw.foppa.dk:28000"}
		assert.Equal(t, "12@qw.foppa.dk:28000", stream.Url())
	})
}
