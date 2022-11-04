package lastscores

import (
	"time"

	"github.com/vikpe/qw-hub-api/pkg/qdemo"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
)

type Entry struct {
	Timestamp    time.Time     `json:"timestamp"`
	Mode         string        `json:"mode"`
	Participants []Participant `json:"players"`
	Map          string        `json:"map"`
}

type Participant struct {
	Name  string
	Frags int
}

func FromLastStatsEntry(entry laststats.Entry) Entry {
	participants := make([]Participant, 0)

	switch entry.Mode {
	case "team":
		fragsPerTeam := make(map[string]int, 0)

		for _, p := range entry.Players {
			if _, ok := fragsPerTeam[p.Team]; !ok {
				fragsPerTeam[p.Team] = 0
			}
			fragsPerTeam[p.Team] += p.Stats.Frags
		}

		for name, frags := range fragsPerTeam {
			participants = append(participants, Participant{
				Name:  name,
				Frags: frags,
			})
		}
	default:
		for _, p := range entry.Players {
			participants = append(participants, Participant{
				Name:  p.Name,
				Frags: p.Stats.Frags,
			})
		}
	}

	timestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", entry.Date)

	return Entry{
		Timestamp:    timestamp,
		Mode:         qdemo.Filename(entry.Demo).Mode(),
		Participants: participants,
		Map:          entry.Map,
	}
}
