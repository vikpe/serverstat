package convert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
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
		Settings: qsettings.Settings{"map": "dm2", "*gamedir": "qw", "status": "Standby", "maxclients": "8", "maxspectators": "4"},
		ExtraInfo: struct {
			QtvStream qtvstream.QtvStream
			Geo       geo.Info
		}{},
	}

	expect := mvdsv.Mvdsv{
		Address: genericServer.Address,
		Mode:    qmode.Mode("ffa"),
		Title:   "ffa [dm2]",
		Status:  "Standby",
		Time: qtime.Time{
			Elapsed:   0,
			Total:     0,
			Remaining: 0,
		},
		Players: []qclient.Client{playerClient},
		PlayerSlots: slots.Slots{
			Used:  1,
			Total: 8,
			Free:  7,
		},
		SpectatorNames: []qstring.QuakeString{spectatorClient.Name},
		SpectatorSlots: slots.Slots{
			Used:  1,
			Total: 4,
			Free:  3,
		},
		Settings:  genericServer.Settings,
		QtvStream: genericServer.ExtraInfo.QtvStream,
		Teams:     []qteam.Team{},
		Geo: geo.Info{
			CC:      "",
			Country: "",
			Region:  "",
		},
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
			Geo       geo.Info
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
			Geo       geo.Info
		}{},
	}

	expect := qwfwd.Qwfwd{
		Address:     genericServer.Address,
		ClientNames: []qstring.QuakeString{client.Name},
		Settings:    genericServer.Settings,
	}

	assert.Equal(t, expect, convert.ToQwfwd(genericServer))
}
