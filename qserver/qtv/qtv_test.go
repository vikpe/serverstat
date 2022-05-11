package qtv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestParse(t *testing.T) {
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

	assert.Equal(t, expect, qtv.Parse(genericServer))
}
