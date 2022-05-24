package mvdsv

import (
	"sort"

	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qstatus"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
	"github.com/vikpe/serverstat/qserver/qtitle"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

const Name = "mvdsv"
const VersionPrefix = Name

type Mvdsv struct {
	Address        string
	Players        []qclient.Client
	SpectatorNames []qstring.QuakeString
	Settings       qsettings.Settings
	QtvStream      qtvstream.QtvStream
	Geo            geo.Info
}

type MvdsvExport struct {
	Type           string
	Address        string
	Mode           qmode.Mode
	Title          string
	Status         string
	Time           qtime.Time
	PlayerSlots    slots.Slots
	Players        []qclient.Client
	Teams          []qteam.Team
	SpectatorSlots slots.Slots
	SpectatorNames []qstring.QuakeString
	Settings       qsettings.Settings
	QtvStream      qtvstream.QtvStream
	Geo            geo.Info
}

func (server Mvdsv) Mode() qmode.Mode {
	return qmode.Parse(server.Settings)
}

func (server Mvdsv) Status() string {
	return qstatus.Parse(server.Settings.Get("status", ""))
}

func (server Mvdsv) PlayerSlots() slots.Slots {
	return slots.New(server.Settings.GetInt("maxclients", 0), len(server.Players))
}

func (server Mvdsv) SpectatorSlots() slots.Slots {
	return slots.New(server.Settings.GetInt("maxspectators", 0), len(server.SpectatorNames))
}

func (server Mvdsv) Time() qtime.Time {
	timelimit := server.Settings.GetInt("timelimit", 0)
	status := server.Settings.Get("status", "")
	return qtime.Parse(timelimit, status)
}

func (server Mvdsv) Teams() []qteam.Team {
	if server.Settings.GetInt("teamplay", 0) > 0 {
		return qteam.FromPlayers(server.Players)
	}

	return make([]qteam.Team, 0)
}

func Export(server Mvdsv) MvdsvExport {
	sort.Slice(server.Players, func(i, j int) bool {
		return server.Players[i].Frags > server.Players[j].Frags
	})

	teams := server.Teams()

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name.ToPlainString() < teams[j].Name.ToPlainString()
	})

	return MvdsvExport{
		Type:           Name,
		Address:        server.Address,
		Mode:           server.Mode(),
		Title:          qtitle.New(server.Settings, server.Players),
		Status:         server.Status(),
		Time:           server.Time(),
		PlayerSlots:    server.PlayerSlots(),
		Players:        server.Players,
		Teams:          teams,
		SpectatorSlots: server.SpectatorSlots(),
		SpectatorNames: server.SpectatorNames,
		Settings:       server.Settings,
		QtvStream:      server.QtvStream,
		Geo:            server.Geo,
	}
}

func GetQtvUsers(address string) ([]qstring.QuakeString, error) {
	return qtvusers.ParseResponse(
		udpclient.New().SendCommand(address, qtvusers.Command),
	)
}

func GetQtvStream(address string) (qtvstream.QtvStream, error) {
	stream, err := status32.ParseResponse(
		udpclient.New().SendCommand(address, status32.Command),
	)

	if err == nil && stream.NumSpectators > 0 {
		spectatorNames, err := GetQtvUsers(address)

		if err == nil {
			stream.SpectatorNames = spectatorNames
		}
	} else {
		stream.SpectatorNames = make([]qstring.QuakeString, 0)
	}

	return stream, err
}
