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

func (q QtvStream) MarshalJSON() ([]byte, error) {
	if "" == q.Url {
		return json.Marshal("")
	} else {
		type QtvStreamJson QtvStream

		return json.Marshal(QtvStreamJson{
			Title:      q.Title,
			Url:        q.Url,
			Clients:    q.Clients,
			NumClients: q.NumClients,
		})
	}
}
