package qtvusers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/qclient"
)

/*
func ParseResponse(responseBody []byte, err error) ([]qclient.Client, error) {
	if err != nil {
		return make([]qclient.Client, 0), err
	} else {
		return ParseResponseBody(responseBody)
	}
}
*/

func TestParseResponse(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		result, err := qtvusers.ParseResponse([]byte{}, errors.New("some error"))
		assert.Equal(t, []qclient.Client{}, result)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("no error", func(t *testing.T) {
		result, err := qtvusers.ParseResponse([]byte(`12 "XantoM"`), nil)
		expect := []qclient.Client{
			{
				Name:    "XantoM",
				NameRaw: []int32{88, 97, 110, 116, 111, 77},
			},
		}
		assert.Equal(t, expect, result)
		assert.Equal(t, err, nil)
	})
}

func TestParseResponseBody(t *testing.T) {
	t.Run("empty response body", func(t *testing.T) {
		result, _ := qtvusers.ParseResponseBody([]byte("\n"))
		assert.Equal(t, []qclient.Client{}, result)
	})

	t.Run("invalid response body", func(t *testing.T) {
		result, _ := qtvusers.ParseResponseBody([]byte("fooooo"))
		assert.Equal(t, []qclient.Client{}, result)
	})

	t.Run("non-empty response body", func(t *testing.T) {
		responseBody := []byte(`12 "XantoM" "player"`)

		result, _ := qtvusers.ParseResponseBody(responseBody)
		expect := []qclient.Client{
			{
				Name:    "XantoM",
				NameRaw: []int32{88, 97, 110, 116, 111, 77},
			},
			{
				Name:    "player",
				NameRaw: []int32{112, 108, 97, 121, 101, 114},
			},
		}

		assert.Equal(t, expect, result)
	})
}
