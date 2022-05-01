package qtvstream

import "encoding/json"

type QtvStream struct {
	Title          string
	Url            string
	SpectatorNames []string
	NumSpectators  uint8
}

func New() QtvStream {
	return QtvStream{
		Title:         "",
		Url:           "",
		NumSpectators: 0,
	}
}

func (node *QtvStream) MarshalJSON() ([]byte, error) {
	if "" == node.Url {
		return json.Marshal(nil)
	} else {
		return json.Marshal(QtvStream{
			Title:          node.Title,
			Url:            node.Url,
			SpectatorNames: node.SpectatorNames,
			NumSpectators:  node.NumSpectators,
		})
	}
}
