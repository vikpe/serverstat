package mvdsv

import (
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/udpclient"
)

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
