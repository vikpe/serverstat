package mvdsv

import (
	"github.com/vikpe/serverstat/qserver/mvdsv/commands"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func GetQtvUsers(address string) ([]qclient.Client, error) {
	return commands.QtvUsers(address)
}

func GetQtvStream(address string) (qtvstream.QtvStream, error) {
	return commands.QtvStream(address)
}
