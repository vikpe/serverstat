package serverstat_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udphelper"
)

func TestClient_GetInfo(t *testing.T) {
	client := serverstat.NewClient()

	t.Run("invalid server address", func(t *testing.T) {
		addresses := []string{"", ":", "a:", ":a", "a:a", "qw.foppa.dk::27501", "qw.foppa.dk::27501:"}

		for _, address := range addresses {
			t.Run(address, func(t *testing.T) {
				server, err := client.GetInfo(address)
				assert.Equal(t, qserver.GenericServer{}, server)
				assert.ErrorContains(t, err, "invalid server address")
			})
		}
	})

	t.Run("response error", func(t *testing.T) {
		server, err := client.GetInfo("foo:666")
		assert.Equal(t, qserver.GenericServer{}, server)
		assert.NotNil(t, err)
	})

	t.Run("successful responses", func(t *testing.T) {
		go func() {
			responseHeader := string([]byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'})
			responseBody := `\hostname\troopers.fi:28501\maxfps\77\pm_ktjump\1\*version\MVDSV 0.35-dev
66 2 4 38 "NL" "" 13 13 "red"
65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc"`
			udphelper.New(":8001").Respond([]byte((responseHeader + responseBody)))
		}()
		time.Sleep(10 * time.Millisecond)

		server, err := client.GetInfo(":8001")
		expectedServer := qserver.GenericServer{
			Version: qversion.New("MVDSV 0.35-dev"),
			Address: ":8001",
			Clients: []qclient.Client{
				{
					Id:     66,
					Name:   qstring.New("NL"),
					Team:   qstring.New("red"),
					Skin:   "",
					Colors: [2]uint8{13, 13},
					Frags:  2,
					Ping:   38,
					Time:   4,
				},
				{
					Id:     65,
					Name:   qstring.New("[ServeMe]"),
					Team:   qstring.New("lqwc"),
					Skin:   "",
					Colors: [2]uint8{12, 11},
					Frags:  -9999,
					Ping:   -666,
					Time:   16,
				},
			},
			Settings: map[string]string{
				"*version":        "MVDSV 0.35-dev",
				"hostname":        "troopers.fi:28501\u0087",
				"hostname_parsed": "troopers.fi:28501",
				"maxfps":          "77",
				"pm_ktjump":       "1",
			},
			Geo: geo.Location{},
			ExtraInfo: struct {
				QtvStream qtvstream.QtvStream `json:"qtv_stream"`
			}{
				QtvStream: qtvstream.QtvStream{
					SpectatorNames: make([]string, 0),
				},
			},
		}
		assert.Equal(t, expectedServer, server)
		assert.Nil(t, err)
	})
}

func TestClient_GetInfoFromMany(t *testing.T) {
	client := serverstat.NewClient()

	t.Run("no valid responses", func(t *testing.T) {
		server := client.GetInfoFromMany([]string{"foo:666", "bar:666"})
		assert.Equal(t, []qserver.GenericServer{}, server)
	})

	t.Run("mixed responses", func(t *testing.T) {
		responseHeader := string([]byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'})

		go func() {
			responseBody := `\maxfps\77\*version\MVDSV 0.35-dev
65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc"`
			udphelper.New(":7001").Respond([]byte((responseHeader + responseBody)))
		}()

		go func() {
			responseBody := `\maxfps\77\*version\MVDSV 0.67`
			udphelper.New(":7002").Respond([]byte((responseHeader + responseBody)))
		}()

		go func() {
			responseBody := `\maxfps\77\*version\qwfwd 0.1`
			udphelper.New(":7003").Respond([]byte((responseHeader + responseBody)))
		}()
		time.Sleep(10 * time.Millisecond)

		servers := client.GetInfoFromMany([]string{":7003", ":7001", ":7002", "foo:666"})

		server1 := qserver.GenericServer{
			Version: qversion.New("MVDSV 0.35-dev"),
			Address: ":7001",
			Clients: []qclient.Client{
				{
					Id:     65,
					Name:   qstring.New("[ServeMe]"),
					Team:   qstring.New("lqwc"),
					Skin:   "",
					Colors: [2]uint8{12, 11},
					Frags:  -9999,
					Ping:   -666,
					Time:   16,
				},
			},
			Settings: map[string]string{
				"*version": "MVDSV 0.35-dev",
				"maxfps":   "77",
			},
			Geo: geo.Location{},
			ExtraInfo: struct {
				QtvStream qtvstream.QtvStream `json:"qtv_stream"`
			}{
				QtvStream: qtvstream.QtvStream{
					SpectatorNames: make([]string, 0),
				},
			},
		}

		server2 := qserver.GenericServer{
			Version: qversion.New("MVDSV 0.67"),
			Address: ":7002",
			Clients: []qclient.Client{},
			Settings: map[string]string{
				"*version": "MVDSV 0.67",
				"maxfps":   "77",
			},
			Geo: geo.Location{},
			ExtraInfo: struct {
				QtvStream qtvstream.QtvStream `json:"qtv_stream"`
			}{
				QtvStream: qtvstream.QtvStream{
					SpectatorNames: make([]string, 0),
				},
			},
		}

		server3 := qserver.GenericServer{
			Version: qversion.New("qwfwd 0.1"),
			Address: ":7003",
			Clients: []qclient.Client{},
			Settings: map[string]string{
				"*version": "qwfwd 0.1",
				"maxfps":   "77",
			},
			Geo: geo.Location{},
			ExtraInfo: struct {
				QtvStream qtvstream.QtvStream `json:"qtv_stream"`
			}{
				QtvStream: qtvstream.QtvStream{
					SpectatorNames: make([]string, 0),
				},
			},
		}

		assert.Equal(t, []qserver.GenericServer{server1, server2, server3}, servers)
	})
}
