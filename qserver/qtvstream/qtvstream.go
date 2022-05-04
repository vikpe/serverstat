package qtvstream

import (
	"encoding/json"

	"github.com/vikpe/serverstat/qserver/qclient"
)

type QtvStream struct {
	Title      string
	Url        string
	Clients    []qclient.Client
	NumClients uint8
}

func New() QtvStream {
	return QtvStream{
		Title:      "",
		Url:        "",
		Clients:    make([]qclient.Client, 0),
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
