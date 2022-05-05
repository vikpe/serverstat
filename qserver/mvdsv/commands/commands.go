package commands

import (
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func QtvUsers(address string) ([]qclient.Client, error) {
	return qtvusers.New(address)
}

func QtvStream(address string) (qtvstream.QtvStream, error) {
	stream, err := status32.New(address)

	if err != nil && stream.NumClients > 0 {
		clients, err := qtvusers.New(address)

		if err == nil {
			stream.Clients = clients
		}
	}

	return stream, err
}
