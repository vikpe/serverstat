package status23_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/commands/status23"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func TestParseResponseBody(t *testing.T) {
	responseBody := []byte(`maxfps\77\pm_ktjump\1\*version\MVDSV 0.35-dev
66 2 4 38 "NL" "" 13 13 "red"
65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc"
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
			IsBot:   true,
		},
	}

	settings, clients := status23.ParseResponseBody(responseBody)
	assert.Equal(t, expectSettings, settings)
	assert.Equal(t, expectClients, clients)
}
