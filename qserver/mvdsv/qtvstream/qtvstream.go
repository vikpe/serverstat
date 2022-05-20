package qtvstream

import (
	"fmt"

	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/serverstat/qutil"
)

type QtvStream struct {
	Title          string
	Id             int
	Address        string
	SpectatorNames []qstring.QuakeString
	NumSpectators  int
}

func (stream QtvStream) Url() string {
	if "" != stream.Address {
		return fmt.Sprintf("%d@%s", stream.Id, stream.Address)
	} else {
		return ""
	}
}

func (stream QtvStream) MarshalJSON() ([]byte, error) {
	return qutil.MarshalNoEscapeHtml(Export(stream))
}

type QtvStreamExport struct {
	Title          string
	Url            string
	Id             int
	Address        string
	SpectatorNames []qstring.QuakeString
	NumSpectators  int
}

func Export(q QtvStream) QtvStreamExport {
	return QtvStreamExport{
		Title:          q.Title,
		Url:            q.Url(),
		Id:             q.Id,
		Address:        q.Address,
		SpectatorNames: q.SpectatorNames,
		NumSpectators:  0,
	}
}
