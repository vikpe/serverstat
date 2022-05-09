package serverstat_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udphelper"
)

func TestGetInfoFromMany(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		responseHeader := string([]byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'})

		go func() {
			responseBody := `\maxfps\77\pm_ktjump\1\*version\MVDSV 0.35-dev
65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc"`
			udphelper.New(":7001").Respond([]byte((responseHeader + responseBody)))
		}()

		go func() {
			responseBody := `\maxfps\77\pm_ktjump\1\*version\MVDSV 0.67`
			udphelper.New(":7002").Respond([]byte((responseHeader + responseBody)))
		}()
		time.Sleep(10 * time.Millisecond)

		servers := serverstat.GetInfoFromMany([]string{":7001", ":7002", "foo:666"})

		server1 := qserver.GenericServer{
			Version: qversion.New("MVDSV 0.35-dev"),
			Address: ":7001",
			Clients: []qclient.Client{
				{
					Name:   qstring.New("[ServeMe]"),
					Team:   qstring.New("lqwc"),
					Skin:   "",
					Colors: [2]uint8{12, 11},
					Frags:  -9999,
					Ping:   -666,
					Time:   16,
					IsBot:  true,
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

		server2 := qserver.GenericServer{
			Version: qversion.New("MVDSV 0.67"),
			Address: ":7002",
			Clients: []qclient.Client{},
			Settings: map[string]string{
				"*version":  "MVDSV 0.67",
				"maxfps":    "77",
				"pm_ktjump": "1",
			},
			ExtraInfo: struct {
				QtvStream qtvstream.QtvStream
			}{},
		}

		assert.Contains(t, servers, server1)
		assert.Contains(t, servers, server2)
		assert.Len(t, servers, 2)
	})
}
