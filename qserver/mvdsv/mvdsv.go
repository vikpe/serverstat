package mvdsv

import (
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/udpclient"
)

func GetQtvUsers(address string) ([]qclient.Client, error) {
	return qtvusers.Send(udpclient.New(), address)
}

func GetQtvStream(address string) (qtvstream.QtvStream, error) {
	udpClient := udpclient.New()
	stream, err := status32.Send(udpClient, address)

	if err != nil && stream.NumClients > 0 {
		clients, err := qtvusers.Send(udpClient, address)

		if err == nil {
			stream.Clients = clients
		}
	}

	return stream, err
}
