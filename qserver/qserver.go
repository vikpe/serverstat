package qserver

import (
	"github.com/vikpe/serverstat/qserver/commands/status23"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/udpclient"
)

type GenericServer struct {
	Version   qversion.Version `json:"-"`
	Address   string
	Clients   []qclient.Client
	Settings  map[string]string
	ExtraInfo struct {
		QtvStream qtvstream.QtvStream
	} `json:"-"`
}

func GetInfo(address string) (GenericServer, error) {
	settings, clients, err := status23.ParseResponse(
		udpclient.New().SendCommand(address, status23.Command),
	)

	if err != nil {
		return GenericServer{}, err
	}

	server := GenericServer{
		Address:  address,
		Version:  qversion.New(settings["*version"]),
		Clients:  clients,
		Settings: settings,
	}

	if server.Version.IsMvdsv() {
		stream, err := mvdsv.GetQtvStream(address)

		if err != nil {
			server.ExtraInfo.QtvStream = stream
		}
	}

	return server, nil
}
