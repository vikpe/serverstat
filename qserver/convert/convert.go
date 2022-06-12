package convert

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/ssoroka/slice"
	"github.com/valyala/fastjson"
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
)

func ToMvdsv(server qserver.GenericServer) mvdsv.Mvdsv {
	players := server.Players()
	qclient.SortPlayers(players)

	spectatorNames := clientNames(server.Spectators())
	status := server.Settings.Get("status", "")
	playerSlots := slots.New(server.Settings.GetInt("maxclients", 0), len(players))
	spectatorSlots := slots.New(server.Settings.GetInt("maxspectators", 0), len(spectatorNames))
	timelimit := server.Settings.GetInt("timelimit", 0)

	teams := make([]qteam.Team, 0)

	if server.Settings.GetInt("teamplay", 0) > 0 {
		teams = qteam.FromPlayers(players)
	}

	return mvdsv.Mvdsv{
		Address:        server.Address,
		Mode:           qmode.Parse(server.Settings),
		Title:          qtitle.New(server.Settings, server.Players()),
		Status:         qstatus.Parse(status),
		Time:           qtime.Parse(timelimit, status),
		Players:        players,
		PlayerSlots:    playerSlots,
		Teams:          teams,
		SpectatorNames: spectatorNames,
		SpectatorSlots: spectatorSlots,
		Settings:       server.Settings,
		QtvStream:      server.ExtraInfo.QtvStream,
		Geo:            server.Geo,
	}
}

func ToQtv(server qserver.GenericServer) qtv.Qtv {
	return qtv.Qtv{
		Address:        server.Address,
		SpectatorNames: clientNames(server.Clients),
		Settings:       server.Settings,
		Geo:            server.Geo,
	}
}

func ToQwfwd(server qserver.GenericServer) qwfwd.Qwfwd {
	return qwfwd.Qwfwd{
		Address:     server.Address,
		ClientNames: clientNames(server.Clients),
		Settings:    server.Settings,
		Geo:         server.Geo,
	}
}

func clientNames(clients []qclient.Client) []string {
	if 0 == len(clients) {
		return make([]string, 0)
	}

	return slice.Map[qclient.Client, string](clients, func(client qclient.Client) string {
		return client.Name.ToPlainString()
	})
}

func ToJson(server qserver.GenericServer) string {
	serverToJson := func(v any) []byte {
		jsonBytes, _ := json.Marshal(v)
		return jsonBytes
	}

	var serverJsonBytes []byte

	if server.Version.IsMvdsv() {
		serverJsonBytes = serverToJson(ToMvdsv(server))
	} else if server.Version.IsQtv() {
		serverJsonBytes = serverToJson(ToQtv(server))
	} else if server.Version.IsQwfwd() {
		serverJsonBytes = serverToJson(ToQwfwd(server))
	} else {
		serverJsonBytes = serverToJson(server)
	}

	value := fastjson.MustParseBytes(serverJsonBytes)
	value.Set("type", fastjson.MustParse(fmt.Sprintf(`"%s"`, server.Version.GetType())))
	buff := value.MarshalTo(nil)
	return string(buff)
}
