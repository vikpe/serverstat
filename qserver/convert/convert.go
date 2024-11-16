package convert

import (
	"fmt"
	"net"

	"github.com/goccy/go-json"
	"github.com/valyala/fastjson"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qscore"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
	"github.com/vikpe/serverstat/qserver/qtitle"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"github.com/vikpe/serverstat/qutil"
)

func ToMvdsv(server qserver.GenericServer) mvdsv.Mvdsv {
	players := qclient.SortPlayers(server.Players())
	spectatorNames := qclient.ClientNames(server.Spectators())
	playerSlots := slots.New(server.Settings.GetInt("maxclients", 0), len(players))
	spectatorSlots := slots.New(server.Settings.GetInt("maxspectators", 0), len(spectatorNames))
	timelimit := server.Settings.GetInt("timelimit", 0)
	settingsStatus := server.Settings.Get("status", "")

	teams := make([]qteam.Team, 0)
	if server.Settings.GetInt("teamplay", 0) > 0 {
		teams = qteam.FromPlayers(players)
	}

	mode, submode := qmode.Parse(server.Settings)
	mvdsvServer := mvdsv.Mvdsv{
		Address:        server.Address,
		Mode:           mode,
		Submode:        submode,
		Title:          qtitle.New(server.Settings, server.Players()),
		Status:         qstatus.New(settingsStatus, mode, server.Players(), playerSlots.Free),
		Time:           qtime.New(timelimit, settingsStatus),
		Players:        players,
		PlayerSlots:    playerSlots,
		Teams:          teams,
		SpectatorNames: spectatorNames,
		SpectatorSlots: spectatorSlots,
		Settings:       server.Settings,
		QtvStream:      server.ExtraInfo.QtvStream,
		Geo:            server.Geo,
	}

	// score: idle server
	if analyze.IsIdle(mvdsvServer) {
		mvdsvServer.Score = 0
	} else {
		mvdsvServer.Score = qscore.FromModeAndPlayers(string(mvdsvServer.Mode), mvdsvServer.Players)
	}

	return mvdsvServer
}

func ToQtv(server qserver.GenericServer) qtv.Qtv {
	return qtv.Qtv{
		Address:        server.Address,
		SpectatorNames: qclient.ClientNames(server.Clients),
		Settings:       server.Settings,
		Geo:            server.Geo,
	}
}

func ToQwfwd(server qserver.GenericServer) qwfwd.Qwfwd {
	return qwfwd.Qwfwd{
		Address:     server.Address,
		ClientNames: qclient.ClientNames(server.Clients),
		Settings:    server.Settings,
		Geo:         server.Geo,
	}
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

	host, port, _ := net.SplitHostPort(server.Address)
	value.Set("host", fastjson.MustParse(fmt.Sprintf(`"%s"`, host)))
	value.Set("port", fastjson.MustParse(fmt.Sprintf(`%d`, qutil.StringToInt(port))))

	buff := value.MarshalTo(nil)
	return string(buff)
}
