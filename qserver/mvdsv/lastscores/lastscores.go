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

	participants := make([]string, 0)
	frags := make([]string, 0)

	if len(e.Teams) > 0 {
		for _, team := range e.Teams {
			participants = append(participants, team.Name.ToPlainString())
			frags = append(frags, strconv.Itoa(team.Frags()))
		}
	} else {
		for _, player := range e.Players {
			participants = append(participants, player.Name.ToPlainString())
			frags = append(frags, strconv.Itoa(player.Frags))
		}
	}

	return fmt.Sprintf("%s - %s: %s [%s] %s",
		timestampStr,
		e.Mode,
		strings.Join(participants, " vs "),
		e.Map,
		strings.Join(frags, ":"),
	)
}

func (e Entry) MarshalJSON() ([]byte, error) {
	type entryJson struct {
		Title     string           `json:"title"`
		Timestamp time.Time        `json:"timestamp"`
		Mode      string           `json:"mode"`
		Players   []qclient.Client `json:"players"`
		Teams     []qteam.Team     `json:"teams"`
		Map       string           `json:"map"`
	}

	return json.Marshal(&entryJson{
		Title:     e.String(),
		Timestamp: e.Timestamp,
		Mode:      e.Mode,
		Players:   e.Players,
		Teams:     e.Teams,
		Map:       e.Map,
	})
}
