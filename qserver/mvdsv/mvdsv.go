package mvdsv

import (
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func GetQtvUsers(address string) ([]qclient.Client, error) {
	return qtvusers.SendTo(address)
}

func GetQtvStream(address string) (qtvstream.QtvStream, error) {
	stream, err := status32.SendTo(address)

	if err != nil && stream.NumClients > 0 {
		clients, err := qtvusers.SendTo(address)

		if err == nil {
			stream.Clients = clients
		}
	}

	return stream, err
}
