package mvdsv

import (
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
	return MvdsvExport{
		Type:           Name,
		Address:        server.Address,
		Mode:           server.Mode(),
		Title:          qtitle.New(server.Settings, server.Players),
		Status:         server.Status(),
		Time:           server.Time(),
		PlayerSlots:    server.PlayerSlots(),
		Players:        server.Players,
		Teams:          server.Teams(),
		SpectatorSlots: server.SpectatorSlots(),
		SpectatorNames: server.SpectatorNames,
		Settings:       server.Settings,
		QtvStream:      server.QtvStream,
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
		names, _ := GetQtvUsers(address)
		stream.SpectatorNames = names
	}

	return stream, err
}
