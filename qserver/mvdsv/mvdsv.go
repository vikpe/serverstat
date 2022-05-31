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
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

const Name = "mvdsv"
const VersionPrefix = Name

type Mvdsv struct {
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
	SpectatorNames []string
	Settings       qsettings.Settings
	QtvStream      qtvstream.QtvStream
	Geo            geo.Info
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
