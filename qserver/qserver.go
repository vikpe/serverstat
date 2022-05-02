package qserver

import (
	"encoding/json"
	"strings"
)

type extraInfo struct {
	QtvStream QtvStream
}

func newExtraInfo() extraInfo {
	return extraInfo{
		QtvStream: NewQtvStream(),
	}
}

type GenericServer struct {
	Address    string
	Clients    []Client
	NumClients uint8
	Settings   map[string]string
	ExtraInfo  extraInfo
}

func NewGenericServer() GenericServer {
	return GenericServer{
		Address:    "",
		Clients:    make([]Client, 0),
		NumClients: 0,
		Settings:   make(map[string]string, 0),
		ExtraInfo:  newExtraInfo(),
	}
}

type QtvStream struct {
	Title          string
	Url            string
	SpectatorNames []string
	NumSpectators  uint8
}

func NewQtvStream() QtvStream {
	return QtvStream{
		Title:          "",
		Url:            "",
		SpectatorNames: make([]string, 0),
		NumSpectators:  0,
	}
}

func (node *QtvStream) MarshalJSON() ([]byte, error) {
	if "" == node.Url {
		return json.Marshal("")
	} else {
		return json.Marshal(QtvStream{
			Title:          node.Title,
			Url:            node.Url,
			SpectatorNames: node.SpectatorNames,
			NumSpectators:  node.NumSpectators,
		})
	}
}

type Client struct {
	Name    string
	NameRaw []rune
	Team    string
	TeamRaw []rune
	Skin    string
	Colors  [2]uint8
	Frags   int
	Ping    int
	Time    uint8
	IsBot   bool
}

func IsGameServer(server GenericServer) bool {
	return !(IsQtvServer(server) || IsProxyServer(server))
}

func IsProxyServer(server GenericServer) bool {
	return strings.HasPrefix(server.Settings["*version"], "qwfwd")
}

func IsQtvServer(server GenericServer) bool {
	return strings.HasPrefix(server.Settings["*version"], "QTV")
}