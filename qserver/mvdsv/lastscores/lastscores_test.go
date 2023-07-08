package lastscores_test

import (
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
	"github.com/vikpe/serverstat/qserver/mvdsv/lastscores"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qteam"
)

var playerXantoM = laststats.Player{Name: "XantoM", Team: "f0m", Ping: 25, TopColor: 4, BottomColor: 2, Stats: laststats.PlayerStats{Frags: 6}}
var playerXterm = laststats.Player{Name: "Xterm", Team: "f0m", Ping: 12, TopColor: 4, BottomColor: 2, Stats: laststats.PlayerStats{Frags: 4}}
var playerBps = laststats.Player{Name: "bps", Team: "-s-", Ping: 25, TopColor: 13, BottomColor: 0, Stats: laststats.PlayerStats{Frags: 20}}
var playerCara = laststats.Player{Name: "carapace", Team: "-s-", Ping: 25, TopColor: 13, BottomColor: 0, Stats: laststats.PlayerStats{Frags: 10}}
var clientXantoM = qclient.NewFromLastStatsPlayer(playerXantoM)
var clientXterm = qclient.NewFromLastStatsPlayer(playerXterm)
var clientBps = qclient.NewFromLastStatsPlayer(playerBps)
var clientCara = qclient.NewFromLastStatsPlayer(playerCara)

func TestNewFromLastStatsEntry(t *testing.T) {
	expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")

	t.Run("duel", func(t *testing.T) {
		laststatsEntry := laststats.Entry{
			Date:    "2022-10-31 19:59:46 +0100",
			Map:     "dm6",
			Mode:    "duel",
			Demo:    "duel_blue_vs_red[dm6]20221031-1958.mvd",
			Teams:   []string{"blue", "red"},
			Players: []laststats.Player{playerXantoM, playerXterm},
		}

		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "duel",
			Demo:      "duel_blue_vs_red[dm6]20221031-1958.mvd",
			Teams:     make([]qteam.Team, 0),
			Players:   []qclient.Client{clientXantoM, clientXterm},
			Map:       "dm6",
		}

		assert.Equal(t, expect, lastscores.NewFromLastStatsEntry(laststatsEntry))
	})

	t.Run("teamplay", func(t *testing.T) {
		laststatsEntry := laststats.Entry{
			Date:    "2022-10-31 19:59:46 +0100",
			Map:     "dm6",
			Mode:    "team",
			Demo:    "2on2_blue_vs_red[dm6]20221031-1958.mvd",
			Teams:   []string{"blue", "red"},
			Players: []laststats.Player{playerXantoM, playerXterm, playerBps, playerCara},
		}

		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "2on2",
			Demo:      "2on2_blue_vs_red[dm6]20221031-1958.mvd",
			Teams:     qteam.FromPlayers([]qclient.Client{clientXantoM, clientXterm, clientBps, clientCara}),
			Players:   []qclient.Client{},
			Map:       "dm6",
		}

		assert.Equal(t, expect, lastscores.NewFromLastStatsEntry(laststatsEntry))
	})

	t.Run("ffa", func(t *testing.T) {
		laststatsEntry := laststats.Entry{
			Date:    "2022-10-31 19:59:46 +0100",
			Map:     "dm6",
			Mode:    "ffa",
			Demo:    "ffa_2[dm6]20221031-1958.mvd",
			Teams:   []string{"blue", "red"},
			Players: []laststats.Player{playerXantoM, playerXterm},
		}

		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "ffa",
			Demo:      "ffa_2[dm6]20221031-1958.mvd",
			Teams:     make([]qteam.Team, 0),
			Players:   []qclient.Client{clientXantoM, clientXterm},
			Map:       "dm6",
		}

		assert.Equal(t, expect, lastscores.NewFromLastStatsEntry(laststatsEntry))
	})
}

func TestEntry_String(t *testing.T) {
	expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")

	t.Run("duel", func(t *testing.T) {
		entry := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "duel",
			Teams:     make([]qteam.Team, 0),
			Players:   []qclient.Client{clientXantoM, clientXterm},
			Map:       "dm6",
		}

		assert.Equal(t, "2022-10-31 19:59 - duel: XantoM vs Xterm [dm6] 6:4", entry.String())
	})

	t.Run("teamplay", func(t *testing.T) {
		entry := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "2on2",
			Teams:     qteam.FromPlayers([]qclient.Client{clientXantoM, clientXterm, clientBps, clientCara}),
			Players:   []qclient.Client{},
			Map:       "dm6",
		}

		assert.Equal(t, "2022-10-31 19:59 - 2on2: -s- vs f0m [dm6] 30:10", entry.String())
	})

	t.Run("ffa", func(t *testing.T) {
		entry := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "ffa",
			Teams:     make([]qteam.Team, 0),
			Players:   []qclient.Client{clientXantoM, clientXterm},
			Map:       "dm6",
		}
		assert.Equal(t, "2022-10-31 19:59 - ffa [dm6]", entry.String())
	})
}

func TestEntry_MarshalJSON(t *testing.T) {
	timestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")
	entry := lastscores.Entry{
		Timestamp: timestamp,
		Demo:      "2on2_-s-_f0m[dm6]20221031-1958.mvd",
		Mode:      "2on2",
		Teams:     qteam.FromPlayers([]qclient.Client{clientXantoM, clientXterm, clientBps, clientCara}),
		Players:   []qclient.Client{},
		Map:       "dm6",
	}
	entryAsJson, err := json.Marshal(entry)
	jsonStr := string(entryAsJson)
	assert.Contains(t, jsonStr, `"title":"2022-10-31 19:59 - 2on2: -s- vs f0m [dm6] 30:10"`)
	assert.Contains(t, jsonStr, `"demo":"2on2_-s-_f0m[dm6]20221031-1958.mvd"`)
	assert.Contains(t, jsonStr, `"timestamp":"2022-10-31 19:59"`)
	assert.Contains(t, jsonStr, `"mode":"2on2"`)
	assert.Contains(t, jsonStr, `"participants":"-s- vs f0m`)
	assert.Contains(t, jsonStr, `"map":"dm6"`)
	assert.Contains(t, jsonStr, `"scores":"30:10"`)
	assert.Nil(t, err)
}
