package qserver

import (
	"encoding/json"
	"strings"
)

type GenericServer struct {
	Version    Version `json:"-"`
	Address    string
	Clients    []Client
	NumClients uint8
	Settings   map[string]string
	ExtraInfo  extraInfo `json:"-"`
}

func NewGenericServer() GenericServer {
	return GenericServer{
		Version:    "",
		Address:    "",
		Clients:    make([]Client, 0),
		NumClients: 0,
		Settings:   make(map[string]string, 0),
		ExtraInfo:  newExtraInfo(),
	}
}

type extraInfo struct {
	QtvStream QtvStream
}

func newExtraInfo() extraInfo {
	return extraInfo{
		QtvStream: NewQtvStream(),
	}
}

type QtvStream struct {
	Title      string
	Url        string
	Clients    []Client
	NumClients uint8
}

func NewQtvStream() QtvStream {
	return QtvStream{
		Title:      "",
		Url:        "",
		Clients:    make([]Client, 0),
		NumClients: 0,
	}
}

func (node *QtvStream) MarshalJSON() ([]byte, error) {
	if "" == node.Url {
		return json.Marshal("")
	} else {
		return json.Marshal(QtvStream{
			Title:      node.Title,
			Url:        node.Url,
			Clients:    node.Clients,
			NumClients: node.NumClients,
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

type ServerType string

const (
	TypeFte     ServerType = "fte"
	TypeMvdsv   ServerType = "mvdsv"
	TypeProxy   ServerType = "qwfwd"
	TypeQtv     ServerType = "qtv"
	TypeUnknown ServerType = "unknown"
)

type Version string

func (sv Version) IsMvdsv() bool {
	return sv.IsType(TypeMvdsv)
}

func (sv Version) IsFte() bool {
	return sv.IsType(TypeFte)
}

func (sv Version) IsProxy() bool {
	return sv.IsType(TypeProxy)
}

func (sv Version) IsQtv() bool {
	return sv.IsType(TypeQtv)
}

func (sv Version) IsGameServer() bool {
	return sv.IsMvdsv() || sv.IsFte()
}

func (sv Version) IsType(serverType ServerType) bool {
	return strings.Contains(
		strings.ToLower(string(sv)),
		strings.ToLower(string(serverType)),
	)
}

func (sv Version) GetType() ServerType {
	if sv.IsProxy() {
		return TypeProxy
	} else if sv.IsMvdsv() {
		return TypeMvdsv
	} else if sv.IsFte() {
		return TypeFte
	} else if sv.IsQtv() {
		return TypeQtv
	} else {
		return TypeUnknown
	}
}
