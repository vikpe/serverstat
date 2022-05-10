package qtvusers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestParseResponse(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		result, err := qtvusers.ParseResponse([]byte{}, errors.New("some error"))
		assert.Equal(t, []qstring.QuakeString{}, result)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("no error", func(t *testing.T) {
		t.Run("invalid response body", func(t *testing.T) {
			t.Run("empty", func(t *testing.T) {
				result, err := qtvusers.ParseResponse([]byte("\n"), nil)
				assert.Equal(t, []qstring.QuakeString{}, result)
				assert.ErrorContains(t, err, "invalid response body")
			})

			t.Run("not containing quotes", func(t *testing.T) {
				result, err := qtvusers.ParseResponse([]byte("fooooo"), nil)
				assert.Equal(t, []qstring.QuakeString{}, result)
				assert.ErrorContains(t, err, "invalid response body")
			})
		})

		t.Run("valid response body", func(t *testing.T) {
			responseBody := []byte(`12 "XantoM" "player"`)

			result, err := qtvusers.ParseResponse(responseBody, nil)
			expect := []qstring.QuakeString{
				qstring.New("XantoM"),
				qstring.New("player"),
			}

			assert.Equal(t, expect, result)
			assert.Nil(t, err)
		})
	})
}
