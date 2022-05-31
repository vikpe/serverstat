package qtvstream

import (
	"github.com/vikpe/serverstat/qtext/qstring"
)

type QtvStream struct {
	Title          string
	Url            string
	Id             int
	Address        string
	SpectatorNames []qstring.QuakeString
	SpectatorCount int
}
