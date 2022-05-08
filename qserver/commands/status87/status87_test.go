package status87_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/commands/status87"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func TestParseResponse(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		settings, clients, err := status87.ParseResponse([]byte{}, errors.New("some error"))
		assert.Equal(t, map[string]string{}, settings)
		assert.Equal(t, []qclient.Client{}, clients)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("no error", func(t *testing.T) {
		responseBody := []byte(`maxfps\77\pm_ktjump\1\*version\MVDSV 0.35-dev
66 2 4 38 "NL" "" 13 13 "red" "SE"
65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc" ""
`)
		expectSettings := map[string]string{
			"*version":  "MVDSV 0.35-dev",
			"maxfps":    "77",
			"pm_ktjump": "1",
		}
		expectClients := []qclient.Client{
			{
				Name:    "NL",
				NameRaw: []rune("NL"),
				Team:    "red",
				TeamRaw: []rune("red"),
				Skin:    "",
				Colors:  [2]uint8{13, 13},
				Frags:   2,
				Ping:    38,
				Time:    4,
				CC:      "SE",
				IsBot:   false,
			},
			{
				Name:    "[ServeMe]",
				NameRaw: []rune("[ServeMe]"),
				Team:    "lqwc",
				TeamRaw: []rune("lqwc"),
				Skin:    "",
				Colors:  [2]uint8{12, 11},
				Frags:   -9999,
				Ping:    -666,
				Time:    16,
				CC:      "",
				IsBot:   true,
			},
		}

		settings, clients, err := status87.ParseResponse(responseBody, nil)
		assert.Equal(t, expectSettings, settings)
		assert.Equal(t, expectClients, clients)
		assert.Nil(t, err)
	})
}
