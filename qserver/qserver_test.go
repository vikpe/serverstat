package qserver_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/udphelper"
)

func TestGetInfo(t *testing.T) {
	t.Run("Response error", func(t *testing.T) {
		server, err := qserver.GetInfo("foo:666")
		assert.Equal(t, qserver.GenericServer{}, server)
		assert.NotNil(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		go func() {
			responseHeader := string([]byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'})
			responseBody := `\maxfps\77\pm_ktjump\1\*version\MVDSV 0.35-dev
66 2 4 38 "NL" "" 13 13 "red"
65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc"`
			udphelper.New(":8001").Respond([]byte((responseHeader + responseBody)))
		}()
		time.Sleep(10 * time.Millisecond)

		server, err := qserver.GetInfo(":8001")
		expectedServer := qserver.GenericServer{
			Version: qversion.New("MVDSV 0.35-dev"),
			Address: ":8001",
			Clients: []qclient.Client{
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
			},
			Settings: map[string]string{
				"*version":  "MVDSV 0.35-dev",
				"maxfps":    "77",
				"pm_ktjump": "1",
			},
			ExtraInfo: struct {
				QtvStream qtvstream.QtvStream
			}{},
		}
		assert.Equal(t, expectedServer, server)
		assert.Nil(t, err)
	})
}
