package qtvstream

import (
	"github.com/vikpe/serverstat/qtext/qstring"
)

type QtvStream struct {
	Title          string                `json:"title"`
	Url            string                `json:"url"`
	ID             int                   `json:"id"`
	Address        string                `json:"address"`
	SpectatorNames []qstring.QuakeString `json:"spectator_names"`
	SpectatorCount int                   `json:"spectator_count"`
}
