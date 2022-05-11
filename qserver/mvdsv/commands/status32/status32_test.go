package status32_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestParseResponse(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		result, err := status32.ParseResponse([]byte(""), errors.New("some error"))
		assert.Equal(t, qtvstream.QtvStream{}, result)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("empty response body", func(t *testing.T) {
		result, err := status32.ParseResponse([]byte(""), nil)
		assert.Equal(t, qtvstream.QtvStream{}, result)
		assert.ErrorContains(t, err, "unable to parse response")
	})

	t.Run("invalid qtv configuration", func(t *testing.T) {
		result, err := status32.ParseResponse([]byte(`1 "qw.foppa.dk - qtv (3)" "" 2`), nil)
		assert.Equal(t, qtvstream.QtvStream{}, result)
		assert.ErrorContains(t, err, "invalid QTV configuration")
	})

	t.Run("valid response body", func(t *testing.T) {
		responseBody := []byte(`1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 4`)

		result, err := status32.ParseResponse(responseBody, nil)
		expect := qtvstream.QtvStream{
			Title:          "qw.foppa.dk - qtv (3)",
			Url:            "3@qw.foppa.dk:28000",
			NumSpectators:  4,
			SpectatorNames: []qstring.QuakeString{},
		}
		assert.Equal(t, expect, result)
		assert.Nil(t, err)
	})
}
