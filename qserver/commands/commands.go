package commands

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/commands/status23"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func GetInfo(address string) (qserver.GenericServer, error) {
	server, err := status23.New(address)

	if server.Version.IsMvdsv() {
		server.ExtraInfo.QtvStream, _ = mvdsv.GetQtvStream(address)
	}

	return server, err
}
