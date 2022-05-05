package qserver

import (
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qversion"
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
