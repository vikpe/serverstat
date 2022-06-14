package mvdsv

import (
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
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

const Name = "mvdsv"
const VersionPrefix = Name

type Mvdsv struct {
	Address        string              `json:"address"`
	Mode           qmode.Mode          `json:"mode"`
	Title          string              `json:"title"`
	Status         qstatus.Status      `json:"status"`
	Time           qtime.Time          `json:"time"`
	PlayerSlots    slots.Slots         `json:"player_slots"`
	Players        []qclient.Client    `json:"players"`
	Teams          []qteam.Team        `json:"teams"`
	SpectatorSlots slots.Slots         `json:"spectator_slots"`
	SpectatorNames []string            `json:"spectator_names"`
	Settings       qsettings.Settings  `json:"settings"`
	QtvStream      qtvstream.QtvStream `json:"qtv_stream"`
	Geo            geo.Info            `json:"geo"`
	Score          int                 `json:"score"`
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

	if err == nil && stream.SpectatorCount > 0 {
		spectatorNames, err := GetQtvUsers(address)

		if err == nil {
			stream.SpectatorNames = spectatorNames
		}
	} else {
		stream.SpectatorNames = make([]qstring.QuakeString, 0)
	}

	stream.SpectatorCount = len(stream.SpectatorNames)

	return stream, err
}
