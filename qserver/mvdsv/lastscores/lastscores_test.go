package lastscores_test

import (
	"testing"
	"time"

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
	t.Run("duel", func(t *testing.T) {
		laststatsEntry := laststats.Entry{
			Date:    "2022-10-31 19:59:46 +0100",
			Map:     "dm6",
			Mode:    "duel",
			Demo:    "duel_blue_vs_red[dm6]20221031-1958.mvd",
			Teams:   []string{"blue", "red"},
			Players: []laststats.Player{playerXantoM, playerXterm},
		}

		expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")
		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "duel",
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

		expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")
		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "2on2",
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

		expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")
		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "ffa",
			Teams:     make([]qteam.Team, 0),
			Players:   []qclient.Client{clientXantoM, clientXterm},
			Map:       "dm6",
		}

		assert.Equal(t, expect, lastscores.NewFromLastStatsEntry(laststatsEntry))
	})
}
