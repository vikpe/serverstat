package bot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qclient/bot"
)

func TestIsBotName(t *testing.T) {
	testCases := map[string]bool{
		"twitch.tv/vikpe": true,
		"[ServeMe]":       true,
		"":                false,
		"XantoM":          false,
	}

	for name, expect := range testCases {
		assert.Equal(t, expect, bot.IsBotName(name))
	}
}

func TestIsBotPing(t *testing.T) {
	testCases := map[int]bool{
		10:   true,
		-550: true,
		12:   false,
		255:  false,
	}

	for ping, expect := range testCases {
		assert.Equal(t, expect, bot.IsBotPing(ping))
	}
}
