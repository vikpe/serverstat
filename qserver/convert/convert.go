package convert

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func ToMvdsv(genericServer qserver.GenericServer) mvdsv.Mvdsv {
	spectatorNames := make([]qstring.QuakeString, 0)
	players := make([]qclient.Client, 0)

	for _, client := range genericServer.Clients {
		if client.IsSpectator() {
			spectatorNames = append(spectatorNames, client.Name)
		} else {
			players = append(players, client)
		}
	}

	return mvdsv.Mvdsv{
		Address:        genericServer.Address,
		Mode:           qmode.Parse(genericServer.Settings),
		Status:         qstatus.Parse(genericServer.Settings),
		Players:        players,
		SpectatorNames: spectatorNames,
		Settings:       genericServer.Settings,
		QtvStream:      genericServer.ExtraInfo.QtvStream,
	}
}

func ToQtv(genericServer qserver.GenericServer) qtv.Qtv {
	spectatorNames := make([]qstring.QuakeString, 0)

	for _, client := range genericServer.Clients {
		spectatorNames = append(spectatorNames, client.Name)
	}

	return qtv.Qtv{
		Address:        genericServer.Address,
		SpectatorNames: spectatorNames,
		Settings:       genericServer.Settings,
	}
}

func ToQwfwd(genericServer qserver.GenericServer) qwfwd.Qwfwd {
	clientNames := make([]qstring.QuakeString, 0)

	for _, client := range genericServer.Clients {
		clientNames = append(clientNames, client.Name)
	}

	return qwfwd.Qwfwd{
		Address:     genericServer.Address,
		ClientNames: clientNames,
		Settings:    genericServer.Settings,
	}
}
