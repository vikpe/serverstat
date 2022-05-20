package qtvstream

import (
	"fmt"

	"github.com/vikpe/serverstat/qtext/qstring"
)

type QtvStream struct {
	Title          string
	Id             int
	Address        string
	SpectatorNames []qstring.QuakeString
	NumSpectators  int
}

func (q QtvStream) Url() string {
	if "" != q.Address {
		return fmt.Sprintf("%d@%s", q.Id, q.Address)
	} else {
		return ""
	}
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
