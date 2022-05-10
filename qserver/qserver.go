package qserver

import (
	"github.com/vikpe/serverstat/qserver/commands/status87"
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
	settings, clients, err := status87.ParseResponse(
		udpclient.New().SendCommand(address, status87.Command),
	)

	if err != nil {
		return GenericServer{}, err
	}

	server := GenericServer{
		Address:  address,
		Version:  qversion.New(settings.Get("*version", "")),
		Clients:  clients,
		Settings: settings,
	}

	if server.Version.IsMvdsv() {
		stream, _ := mvdsv.GetQtvStream(address)
		server.ExtraInfo.QtvStream = stream
	}

	return server, nil
}
