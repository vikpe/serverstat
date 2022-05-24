package qtvstream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qtext/qstring"
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

func TestExport(t *testing.T) {
	stream := qtvstream.QtvStream{
		Title:          "foppa qtv",
		Id:             12,
		Address:        "qw.foppa.dk:28000",
		SpectatorNames: []qstring.QuakeString{qstring.New("XantoM")},
		NumSpectators:  1,
	}
	expect := qtvstream.QtvStreamExport{
		Title:          stream.Title,
		Url:            "12@qw.foppa.dk:28000",
		Id:             stream.Id,
		Address:        stream.Address,
		SpectatorNames: stream.SpectatorNames,
		NumSpectators:  1,
	}
	assert.Equal(t, expect, qtvstream.Export(stream))
}
