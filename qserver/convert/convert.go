package convert

import (
	"github.com/ssoroka/slice"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qstatus"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
	"github.com/vikpe/serverstat/qserver/qtitle"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func ToMvdsv(server qserver.GenericServer) mvdsv.Mvdsv {
	spectatorNames := clientNames(server.Spectators())
	status := qstatus.Parse(server.Settings.Get("status", ""))
	playerSlots := slots.New(server.Settings.GetInt("maxclients", 0), len(server.Players()))
	spectatorSlots := slots.New(server.Settings.GetInt("maxspectators", 0), len(spectatorNames))
	timelimit := server.Settings.GetInt("timelimit", 0)

	players := server.Players()
	qclient.SortPlayers(players)

	teams := make([]qteam.Team, 0)

	if server.Settings.GetInt("teamplay", 0) > 0 {
		teams = qteam.FromPlayers(players)
	}

	return mvdsv.Mvdsv{
		Address:        server.Address,
		Mode:           qmode.Parse(server.Settings),
		Title:          qtitle.New(server.Settings, server.Players()),
		Status:         status,
		Time:           qtime.Parse(timelimit, status),
		Players:        players,
		PlayerSlots:    playerSlots,
		Teams:          teams,
		SpectatorNames: spectatorNames,
		SpectatorSlots: spectatorSlots,
		Settings:       server.Settings,
		QtvStream:      server.ExtraInfo.QtvStream,
		Geo:            server.ExtraInfo.Geo,
	}
}

func ToQtv(server qserver.GenericServer) qtv.Qtv {
	return qtv.Qtv{
		Address:        server.Address,
		SpectatorNames: clientNames(server.Clients),
		Settings:       server.Settings,
		Geo:            server.ExtraInfo.Geo,
	}
}

func ToQwfwd(server qserver.GenericServer) qwfwd.Qwfwd {
	return qwfwd.Qwfwd{
		Address:     server.Address,
		ClientNames: clientNames(server.Clients),
		Settings:    server.Settings,
		Geo:         server.ExtraInfo.Geo,
	}
}

func clientNames(clients []qclient.Client) []qstring.QuakeString {
	return slice.Map[qclient.Client, qstring.QuakeString](clients, func(client qclient.Client) qstring.QuakeString {
		return client.Name
	})
}
