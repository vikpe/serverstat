package mvdsv_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udphelper"
)

func TestGetQtvStream(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		stream, err := mvdsv.GetQtvStream("foo:666")
		assert.Equal(t, qtvstream.QtvStream{}, stream)
		assert.NotNil(t, err)
	})

	t.Run("has no clients", func(t *testing.T) {
		go func() {
			qtvHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', ' '}
			qtvBody := []byte(`1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 0`)
			qtvResponse := append(qtvHeader, qtvBody...)

			udphelper.New(":5001").Respond(qtvResponse)
		}()
		time.Sleep(10 * time.Millisecond)

		stream, err := mvdsv.GetQtvStream(":5001")
		expectStream := qtvstream.QtvStream{
			Title:          "qw.foppa.dk - qtv (3)",
			Url:            "3@qw.foppa.dk:28000",
			NumSpectators:  0,
			SpectatorNames: []qstring.QuakeString{},
		}
		assert.Equal(t, expectStream, stream)
		assert.Nil(t, err)
	})

	t.Run("has clients", func(t *testing.T) {
		go func() {
			qtvHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', ' '}
			qtvBody := []byte(`1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 2`)
			qtvResponse := append(qtvHeader, qtvBody...)

			qtvusersHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
			qtvusersBody := []byte(`12 "XantoM" "valla"`)
			qtvusersResponse := append(qtvusersHeader, qtvusersBody...)

			udphelper.New(":5002").Respond(qtvResponse, qtvusersResponse)
		}()
		time.Sleep(10 * time.Millisecond)

		stream, err := mvdsv.GetQtvStream(":5002")
		expectStream := qtvstream.QtvStream{
			Title:         "qw.foppa.dk - qtv (3)",
			Url:           "3@qw.foppa.dk:28000",
			NumSpectators: 2,
			SpectatorNames: []qstring.QuakeString{
				qstring.New("XantoM"),
				qstring.New("valla"),
			},
		}

		assert.Equal(t, expectStream, stream)
		assert.Nil(t, err)
	})
}

func TestParse(t *testing.T) {
	playerClient := qclient.Client{
		Name:   qstring.New("NL"),
		Team:   qstring.New("red"),
		Skin:   "",
		Colors: [2]uint8{13, 13},
		Frags:  2,
		Ping:   38,
		Time:   4,
	}

	spectatorClient := qclient.Client{
		Name:   qstring.New("[ServeMe]"),
		Team:   qstring.New("lqwc"),
		Skin:   "",
		Colors: [2]uint8{12, 11},
		Frags:  -9999,
		Ping:   -666,
		Time:   16,
	}

	genericServer := qserver.GenericServer{
		Address:  "qw.foppa.dk:27501",
		Version:  qversion.Version("mvdsv 0.15"),
		Clients:  []qclient.Client{playerClient, spectatorClient},
		Settings: qsettings.Settings{"map": "dm2", "*gamedir": "qw", "status": "Standby"},
		ExtraInfo: struct {
			QtvStream qtvstream.QtvStream
		}{},
	}

	expect := mvdsv.Mvdsv{
		Address: genericServer.Address,
		Status: qstatus.Status{
			Name: "Standby",
			Duration: qstatus.MatchDuration{
				Elapsed:   0,
				Total:     0,
				Remaining: 0,
			},
		},
		Mode:           qmode.Mode("ffa"),
		Players:        []qclient.Client{playerClient},
		SpectatorNames: []qstring.QuakeString{spectatorClient.Name},
		Settings:       genericServer.Settings,
		QtvStream:      genericServer.ExtraInfo.QtvStream,
	}

	assert.Equal(t, expect, mvdsv.Parse(genericServer))
}
