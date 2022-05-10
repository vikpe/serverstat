package mvdsv

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

type Server struct {
	Address        string
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
		Players:        players,
		SpectatorNames: spectatorNames,
		Settings:       genericServer.Settings,
		QtvStream:      genericServer.ExtraInfo.QtvStream,
	}
}

func GetQtvUsers(address string) ([]qclient.Client, error) {
	return qtvusers.ParseResponse(
		udpclient.New().SendCommand(address, qtvusers.Command),
	)
}

func GetQtvStream(address string) (qtvstream.QtvStream, error) {
	stream, err := status32.ParseResponse(
		udpclient.New().SendCommand(address, status32.Command),
	)

	if err == nil && stream.NumClients > 0 {
		clients, _ := GetQtvUsers(address)
		stream.Clients = clients
	}

	return stream, err
}
