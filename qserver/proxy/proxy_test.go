package proxy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/proxy"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
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

	expect := proxy.Proxy{
		Address:     genericServer.Address,
		ClientNames: []qstring.QuakeString{client.Name},
		Settings:    genericServer.Settings,
	}

	assert.Equal(t, expect, proxy.Parse(genericServer))
}
