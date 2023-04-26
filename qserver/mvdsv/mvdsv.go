package mvdsv

import (
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/lastscores"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
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
	Geo            geo.Location        `json:"geo"`
	Score          int                 `json:"score"`
}

func GetLastScores(address string, limit int) ([]lastscores.Entry, error) {
	stats, err := GetLastStats(address, limit)

	if err != nil {
		return make([]lastscores.Entry, 0), err
	}

	result := make([]lastscores.Entry, 0)
	for _, entry := range stats {
		result = append(result, lastscores.NewFromLastStatsEntry(entry))
	}

	return result, nil
}

func GetLastStats(address string, limit int) ([]laststats.Entry, error) {
	udpClient := udpclient.New()
	udpClient.Config.BufferSize = 100 * 8192
	return laststats.ParseResponseBody(udpClient.SendCommand(address, laststats.GetCommand(limit)))
}

func GetQtvUsers(address string) ([]string, error) {
	return qtvusers.ParseResponse(
		udpclient.New().SendCommand(address, qtvusers.Command),
	)
}

func GetQtvStream(address string) (qtvstream.QtvStream, error) {
	response, err := udpclient.New().SendCommand(address, status32.Command)

	if err != nil {
		return qtvstream.New(), err
	}

	stream, err := status32.ParseResponse(address, response)

	if err == nil && stream.SpectatorCount > 0 {
		spectatorNames, _ := GetQtvUsers(address)
		stream.SpectatorNames = spectatorNames
	}

	stream.SpectatorCount = len(stream.SpectatorNames)

	return stream, err
}
