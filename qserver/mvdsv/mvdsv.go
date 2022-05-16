package mvdsv

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qstatus"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
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

type ClientSlots struct {
	Used  int
	Total int
	Free  int
}

func NewSlots(total int, used int) ClientSlots {
	return ClientSlots{
		Used:  used,
		Total: total,
		Free:  total - used,
	}
}

func (server Mvdsv) Mode() qmode.Mode {
	return qmode.Parse(server.Settings)
}

func (server Mvdsv) Status() string {
	return qstatus.Parse(server.Settings.Get("status", ""))
}

func (server Mvdsv) PlayerSlots() ClientSlots {
	return NewSlots(server.Settings.GetInt("maxclients", 0), len(server.Players))
}

func (server Mvdsv) SpectatorSlots() ClientSlots {
	return NewSlots(server.Settings.GetInt("maxspectators", 0), len(server.SpectatorNames))
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

func (server Mvdsv) Title() string {
	titleParts := make([]string, 0)

	// matchtag
	matchTag := server.Settings.Get("matchtag", "")

	if matchTag != "" {
		titleParts = append(titleParts, fmt.Sprintf("%s / ", matchTag))
	}

	// mode
	mode := server.Mode()
	titleParts = append(titleParts, fmt.Sprintf("%s:", string(mode)))

	// participants
	participants := make([]string, 0)
	isTeamplay := server.Settings.GetInt("teamplay", 0) > 0

	if isTeamplay && !mode.IsCoop() {
		for _, t := range server.Teams() {
			participants = append(participants, t.String())
		}
	} else if !mode.IsFfa() {
		for _, p := range server.Players {
			participants = append(participants, p.Name.ToPlainString())
		}
	}

	if len(participants) > 0 {
		var participantDelimiter string

		if mode.IsCoop() {
			participantDelimiter = ", "
		} else {
			participantDelimiter = " vs "
		}

		titleParts = append(titleParts, strings.Join(participants, participantDelimiter))
	}

	// map
	titleParts = append(titleParts, fmt.Sprintf("[%s]", server.Settings.Get("map", "")))

	return strings.Join(titleParts, " ")
}

func (server Mvdsv) MarshalJSON() ([]byte, error) {
	type mvdsvJson struct {
		Address        string
		Type           string
		Mode           qmode.Mode
		Status         string
		Time           qtime.Time
		Title          string
		PlayerSlots    ClientSlots
		Players        []qclient.Client
		Teams          []qteam.Team
		SpectatorSlots ClientSlots
		SpectatorNames []qstring.QuakeString
		Settings       qsettings.Settings
		QtvStream      qtvstream.QtvStream
	}

	return json.Marshal(mvdsvJson{
		Address:        server.Address,
		Type:           Name,
		Mode:           server.Mode(),
		Title:          server.Title(),
		Status:         server.Status(),
		Time:           server.Time(),
		PlayerSlots:    server.PlayerSlots(),
		Players:        server.Players,
		Teams:          server.Teams(),
		SpectatorSlots: server.SpectatorSlots(),
		SpectatorNames: server.SpectatorNames,
		Settings:       server.Settings,
		QtvStream:      server.QtvStream,
	})
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
