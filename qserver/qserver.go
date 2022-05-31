package qserver

import (
	"github.com/ssoroka/slice"
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

func (server GenericServer) Players() []qclient.Client {
	return slice.Select(server.Clients, func(i int, c qclient.Client) bool {
		return c.IsPlayer()
	})
}

func (server GenericServer) Spectators() []qclient.Client {
	return slice.Select(server.Clients, func(i int, c qclient.Client) bool {
		return c.IsSpectator()
	})
}
