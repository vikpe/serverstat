package qserver

import (
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qversion"
)

type GenericServer struct {
	Address   string
	Version   qversion.Version
	Clients   []qclient.Client
	Settings  qsettings.Settings
	ExtraInfo struct {
		QtvStream qtvstream.QtvStream
		Geo       geo.Info
	}
}

func (server GenericServer) HasPlayers() bool {
	for _, c := range server.Clients {
		if !c.IsSpectator() && !c.IsBot() {
			return true
		}
	}

	return false
}
