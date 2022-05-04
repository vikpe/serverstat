package qserver

import (
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qtvstream"
	"github.com/vikpe/serverstat/qserver/qversion"
)

type GenericServer struct {
	Version    qversion.Version `json:"-"`
	Address    string
	Clients    []qclient.Client
	NumClients uint8
	Settings   map[string]string
	ExtraInfo  extraInfo `json:"-"`
}

func NewGenericServer() GenericServer {
	return GenericServer{
		Version:    qversion.New(""),
		Address:    "",
		Clients:    make([]qclient.Client, 0),
		NumClients: 0,
		Settings:   make(map[string]string, 0),
		ExtraInfo:  newExtraInfo(),
	}
}

type extraInfo struct {
	QtvStream qtvstream.QtvStream
}

func newExtraInfo() extraInfo {
	return extraInfo{
		QtvStream: qtvstream.New(),
	}
}
