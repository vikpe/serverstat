package mvdsv

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

type Server struct {
	Address        string
	Status         qstatus.Status
	Mode           qmode.Mode
	Players        []qclient.Client
	SpectatorNames []qstring.QuakeString
	Settings       qsettings.Settings
	QtvStream      qtvstream.QtvStream
}

func Parse(genericServer qserver.GenericServer) Server {
	spectatorNames := make([]qstring.QuakeString, 0)
	players := make([]qclient.Client, 0)

	for _, client := range genericServer.Clients {
		if client.IsSpectator() {
			spectatorNames = append(spectatorNames, client.Name)
		} else {
			players = append(players, client)
		}
	}

	return Server{
		Address:        genericServer.Address,
		Mode:           qmode.Parse(genericServer.Settings),
		Status:         qstatus.Parse(genericServer.Settings),
		Players:        players,
		SpectatorNames: spectatorNames,
		Settings:       genericServer.Settings,
		QtvStream:      genericServer.ExtraInfo.QtvStream,
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
