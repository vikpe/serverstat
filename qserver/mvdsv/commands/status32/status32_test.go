package status32_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
)

var EmptyStream = qtvstream.QtvStream{
	Title:          "",
	Url:            "",
	ID:             0,
	Address:        "",
	SpectatorNames: make([]string, 0),
	SpectatorCount: 0,
}

func TestParseResponse(t *testing.T) {
	t.Run("empty response body", func(t *testing.T) {
		result, err := status32.ParseResponse("qw.foppa.dk:28501", []byte(""))
		assert.Equal(t, EmptyStream, result)
		assert.ErrorContains(t, err, "unable to parse response")
	})

	t.Run("invalid qtv configuration", func(t *testing.T) {
		t.Run("URL unparsable from title", func(t *testing.T) {
			result, err := status32.ParseResponse("qw.foppa.dk:28501", []byte(`1 "qw.foppa.dk" "" 2`))
			assert.Equal(t, EmptyStream, result)
			assert.ErrorContains(t, err, "invalid QTV configuration")
		})

		t.Run("URL parsable from title", func(t *testing.T) {
			result, err := status32.ParseResponse("qw.foppa.dk:28501", []byte(`1 "qw.foppa.dk - qtv (3)" "" 4`))
			expect := qtvstream.QtvStream{
				Title:          "qw.foppa.dk - qtv (3)",
				Url:            "3@qw.foppa.dk:28000",
				ID:             3,
				Address:        "qw.foppa.dk:28000",
				SpectatorCount: 4,
				SpectatorNames: make([]string, 0),
			}
			assert.Equal(t, expect, result)
			assert.Nil(t, err)
		})
	})

	t.Run("valid response body", func(t *testing.T) {
		result, err := status32.ParseResponse("qw.foppa.dk:28501", []byte(`1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 4`))
		expect := qtvstream.QtvStream{
			Title:          "qw.foppa.dk - qtv (3)",
			Url:            "3@qw.foppa.dk:28000",
			ID:             3,
			Address:        "qw.foppa.dk:28000",
			SpectatorCount: 4,
			SpectatorNames: make([]string, 0),
		}
		assert.Equal(t, expect, result)
		assert.Nil(t, err)
	})
}

func TestStreamNumberFromTitle(t *testing.T) {
	t.Run("unable to parse stream number", func(t *testing.T) {
		t.Run("no number or braces present", func(t *testing.T) {
			number, err := status32.StreamNumberFromTitle("qw.foppa.dk")
			assert.Equal(t, 0, number)
			assert.ErrorContains(t, err, "unable to parse stream number from title")
		})

		t.Run("braces but no number", func(t *testing.T) {
			number, err := status32.StreamNumberFromTitle("qw.foppa.dk - qtv ()")
			assert.Equal(t, 0, number)
			assert.ErrorContains(t, err, "unable to parse stream number from title")
		})
	})

	t.Run("able to parse stream number", func(t *testing.T) {
		t.Run("single digit", func(t *testing.T) {
			number, err := status32.StreamNumberFromTitle("DuelMania FRANCE Qtv (1)")
			assert.Equal(t, 1, number)
			assert.Nil(t, err)
		})

		t.Run("double digit", func(t *testing.T) {
			number, err := status32.StreamNumberFromTitle("qw.foppa.dk - qtv (13)")
			assert.Equal(t, 13, number)
			assert.Nil(t, err)
		})
	})
}
