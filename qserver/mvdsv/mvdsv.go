package mvdsv

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qteam"
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

func (server Mvdsv) Mode() qmode.Mode {
	return qmode.Parse(server.Settings)
}

func (server Mvdsv) Status() qstatus.Status {
	return qstatus.Parse(server.Settings)
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
		Status         qstatus.Status
		Title          string
		Players        []qclient.Client
		Teams          []qteam.Team
		SpectatorNames []qstring.QuakeString
		Settings       qsettings.Settings
		QtvStream      qtvstream.QtvStream
	}

	return json.Marshal(mvdsvJson{
		Address:        server.Address,
		Type:           Name,
		Mode:           server.Mode(),
		Status:         server.Status(),
		Title:          server.Title(),
		Players:        server.Players,
		Teams:          server.Teams(),
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
