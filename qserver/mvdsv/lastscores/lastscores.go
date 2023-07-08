package lastscores

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qteam"
)

type Entry struct {
	Timestamp time.Time
	Demo      string
	Mode      string
	Players   []qclient.Client
	Teams     []qteam.Team
	Map       string
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
		Demo:      entry.Demo,
		Mode:      qdemo.Filename(entry.Demo).Mode(),
		Teams:     teams,
		Players:   players,
		Map:       entry.Map,
	}
}

func (e Entry) String() string {
	timestampStr := e.Timestamp.Format("2006-01-02 15:04")

	if "ffa" == e.Mode {
		return fmt.Sprintf("%s - %s [%s]", timestampStr, e.Mode, e.Map)
	}

	return fmt.Sprintf("%s - %s: %s [%s] %s",
		timestampStr,
		e.Mode,
		e.Participants(),
		e.Map,
		e.Scores(),
	)
}

func (e Entry) Participants() string {
	participants := make([]string, 0)

	if len(e.Teams) > 0 {
		for _, team := range e.Teams {
			participants = append(participants, team.Name.ToPlainString())
		}
	} else {
		for _, player := range e.Players {
			participants = append(participants, player.Name.ToPlainString())
		}
	}

	return strings.Join(participants, " vs ")
}

func (e Entry) Scores() string {
	frags := make([]string, 0)

	if len(e.Teams) > 0 {
		for _, team := range e.Teams {
			frags = append(frags, strconv.Itoa(team.Frags()))
		}
	} else {
		for _, player := range e.Players {
			frags = append(frags, strconv.Itoa(player.Frags))
		}
	}

	return strings.Join(frags, ":")
}

func (e Entry) MarshalJSON() ([]byte, error) {
	type entryJson struct {
		Title        string           `json:"title"`
		Demo         string           `json:"demo"`
		Timestamp    string           `json:"timestamp"`
		TimestampIso time.Time        `json:"timestamp_iso"`
		Mode         string           `json:"mode"`
		Participants string           `json:"participants"`
		Map          string           `json:"map"`
		Scores       string           `json:"scores"`
		Players      []qclient.Client `json:"players"`
		Teams        []qteam.Team     `json:"teams"`
	}

	return json.Marshal(&entryJson{
		Title:        e.String(),
		Demo:         e.Demo,
		Timestamp:    e.Timestamp.Format("2006-01-02 15:04"),
		TimestampIso: e.Timestamp,
		Mode:         e.Mode,
		Participants: e.Participants(),
		Map:          e.Map,
		Scores:       e.Scores(),
		Players:      e.Players,
		Teams:        e.Teams,
	})
}
