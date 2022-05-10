package qtvstream

import (
	"encoding/json"

	"github.com/vikpe/serverstat/qtext/qstring"
)

type QtvStream struct {
	Title          string
	Url            string
	SpectatorNames []qstring.QuakeString
	NumSpectators  uint8
}

func (q QtvStream) MarshalJSON() ([]byte, error) {
	if "" == q.Url {
		return json.Marshal("")
	} else {
		type QtvStreamJson QtvStream

		return json.Marshal(QtvStreamJson{
			Title:          q.Title,
			Url:            q.Url,
			SpectatorNames: q.SpectatorNames,
			NumSpectators:  q.NumSpectators,
		})
	}
}
