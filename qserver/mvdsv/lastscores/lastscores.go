package lastscores

import (
	"time"

	"github.com/vikpe/qw-hub-api/pkg/qdemo"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qteam"
)

type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Mode      string           `json:"mode"`
	Players   []qclient.Client `json:"players"`
	Teams     []qteam.Team     `json:"teams"`
	Map       string           `json:"map"`
}

func NewFromLastStatsEntry(entry laststats.Entry) Entry {
	clients := make([]qclient.Client, 0)
	for _, p := range entry.Players {
		clients = append(clients, qclient.NewFromLastStatsPlayer(p))
	}

	teams := make([]qteam.Team, 0)
	players := make([]qclient.Client, 0)

	if "team" == entry.Mode {
		teams = qteam.FromPlayers(clients)
	} else {
		players = clients
	}

	timestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", entry.Date)

	return Entry{
		Timestamp: timestamp,
		Mode:      qdemo.Filename(entry.Demo).Mode(),
		Teams:     teams,
		Players:   players,
		Map:       entry.Map,
	}
}
