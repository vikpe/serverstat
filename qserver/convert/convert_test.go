package convert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestToMvdsv(t *testing.T) {
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
		Address:        genericServer.Address,
		Players:        []qclient.Client{playerClient},
		SpectatorNames: []qstring.QuakeString{spectatorClient.Name},
		Settings:       genericServer.Settings,
		QtvStream:      genericServer.ExtraInfo.QtvStream,
	}

	assert.Equal(t, expect, convert.ToMvdsv(genericServer))
}

func TestToQtv(t *testing.T) {
	client := qclient.Client{
		Name:   qstring.New("XantoM"),
		Team:   qstring.New(""),
		Skin:   "",
		Colors: [2]uint8{0, 0},
		Frags:  0,
		Ping:   666,
		Time:   4,
	}

	genericServer := qserver.GenericServer{
		Address: "qw.foppa.dk:28000",
		Version: qversion.Version("QTV 1.12-rc1"),
		Clients: []qclient.Client{client},
		Settings: qsettings.Settings{
			"*version":   "QTV 1.12-rc1",
			"hostname":   "qw.foppa.dk - qtv",
			"maxclients": "100",
		},
		ExtraInfo: struct {
			QtvStream qtvstream.QtvStream
		}{},
	}

	expect := qtv.Qtv{
		Address:        genericServer.Address,
		SpectatorNames: []qstring.QuakeString{client.Name},
		Settings:       genericServer.Settings,
	}

	assert.Equal(t, expect, convert.ToQtv(genericServer))
}

func TestToQwfwd(t *testing.T) {
	client := qclient.Client{
		Name:   qstring.New("XantoM"),
		Team:   qstring.New(""),
		Skin:   "",
		Colors: [2]uint8{0, 0},
		Frags:  0,
		Ping:   666,
		Time:   4,
	}

	genericServer := qserver.GenericServer{
		Address: "troopers.fi:28000",
		Version: qversion.Version("qwfwd 1.2"),
		Clients: []qclient.Client{client},
		Settings: qsettings.Settings{
			"*version":   "qwfwd 1.2",
			"hostname":   "troopers.fi QWfwd",
			"maxclients": "128",
		},
		ExtraInfo: struct {
			QtvStream qtvstream.QtvStream
		}{},
	}

	expect := qwfwd.Qwfwd{
		Address:     genericServer.Address,
		ClientNames: []qstring.QuakeString{client.Name},
		Settings:    genericServer.Settings,
	}

	assert.Equal(t, expect, convert.ToQwfwd(genericServer))
}
