package qclient_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func TestNewFromString(t *testing.T) {
	// valid
	expect := qclient.Client{
		Name:    "XantoM",
		NameRaw: []rune("XantoM"),
		Team:    "f0m",
		TeamRaw: []rune("f0m"),
		Skin:    "xantom",
		Colors:  [2]uint8{4, 2},
		Frags:   17,
		Ping:    12,
		Time:    25,
		IsBot:   false,
	}
	clientString := `585 17 25 12 "XantoM" "xantom" 4 2 "f0m"`
	client, _ := qclient.NewFromString(clientString)
	assert.Equal(t, expect, client)

	// missing fields
	client, err := qclient.NewFromString("585 17 25")
	assert.Equal(t, err, errors.New("invalid client column count 3, expects at least 8"))
	assert.Equal(t, client, qclient.Client{})

	// invalid
	client, err = qclient.NewFromString("")
	assert.Equal(t, err, errors.New("EOF"))
	assert.Equal(t, client, qclient.Client{})
}

func TestFromStrings(t *testing.T) {
	clientStrings := []string{
		`63 5 4 25 "Pitbull" "" 4 4 "red"`,
		`66 2 4 38 "NL" "" 13 13 "red"`,
		`65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc"`,
		`67 -9999 122 -68 "\s\Final" "" 2 3 "red"`,
		``,
	}

	expect := []qclient.Client{
		{
			Name:    "Pitbull",
			NameRaw: []rune("Pitbull"),
			Team:    "red",
			TeamRaw: []rune("red"),
			Skin:    "",
			Colors:  [2]uint8{4, 4},
			Frags:   5,
			Ping:    25,
			Time:    4,
			IsBot:   false,
		},
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
		{
			Name:    "Final",
			NameRaw: []rune("Final"),
			Team:    "red",
			TeamRaw: []rune("red"),
			Skin:    "",
			Colors:  [2]uint8{2, 3},
			Frags:   -9999,
			Ping:    -68,
			Time:    122,
			IsBot:   false,
		},
	}

	actual := qclient.NewFromStrings(clientStrings)

	assert.Equal(t, expect, actual)
}

func TestIsBotName(t *testing.T) {
	knownBotNames := []string{
		"[ServeMe]",
		"twitch.tv/vikpe",
	}

	for _, name := range knownBotNames {
		assert.True(t, qclient.IsBotName(name), name)
	}

	assert.False(t, qclient.IsBotName(""))
	assert.False(t, qclient.IsBotName("XantoM"))
}

func TestIsBotPing(t *testing.T) {
	assert.True(t, qclient.IsBotPing(10))

	assert.False(t, qclient.IsBotPing(0))
	assert.False(t, qclient.IsBotPing(38))
}
