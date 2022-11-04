package lastscores_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
	"github.com/vikpe/serverstat/qserver/mvdsv/lastscores"
)

func TestFromLastStatsEntry(t *testing.T) {
	t.Run("duel", func(t *testing.T) {
		laststatsEntry := laststats.Entry{
			Date:  "2022-10-31 19:59:46 +0100",
			Map:   "dm6",
			Mode:  "duel",
			Demo:  "duel_blue_vs_red[dm6]20221031-1958.mvd",
			Teams: []string{"blue", "red"},
			Players: []laststats.Player{
				{Name: "XantoM", Team: "blue", Stats: laststats.PlayerStats{Frags: 6}},
				{Name: "Xterm", Team: "blue", Stats: laststats.PlayerStats{Frags: 4}},
			},
		}

		expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")
		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "duel",
			Participants: []lastscores.Participant{
				{Name: "XantoM", Frags: 6},
				{Name: "Xterm", Frags: 4},
			},
			Map: "dm6",
		}

		assert.Equal(t, expect, lastscores.FromLastStatsEntry(laststatsEntry))
	})

	t.Run("teamplay", func(t *testing.T) {
		laststatsEntry := laststats.Entry{
			Date:  "2022-10-31 19:59:46 +0100",
			Map:   "dm6",
			Mode:  "team",
			Demo:  "2on2_blue_vs_red[dm6]20221031-1958.mvd",
			Teams: []string{"blue", "red"},
			Players: []laststats.Player{
				{Name: "XantoM", Team: "blue", Stats: laststats.PlayerStats{Frags: 6}},
				{Name: "Xterm", Team: "blue", Stats: laststats.PlayerStats{Frags: 4}},
				{Name: "bps", Team: "red", Stats: laststats.PlayerStats{Frags: 20}},
				{Name: "carapace", Team: "red", Stats: laststats.PlayerStats{Frags: 10}},
			},
		}

		expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")
		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "2on2",
			Participants: []lastscores.Participant{
				{Name: "blue", Frags: 10},
				{Name: "red", Frags: 30},
			},
			Map: "dm6",
		}

		assert.Equal(t, expect, lastscores.FromLastStatsEntry(laststatsEntry))
	})

	t.Run("ffa", func(t *testing.T) {
		laststatsEntry := laststats.Entry{
			Date:  "2022-10-31 19:59:46 +0100",
			Map:   "dm6",
			Mode:  "ffa",
			Demo:  "ffa_2[dm6]20221031-1958.mvd",
			Teams: []string{"blue", "red"},
			Players: []laststats.Player{
				{Name: "XantoM", Team: "blue", Stats: laststats.PlayerStats{Frags: 6}},
				{Name: "Xterm", Team: "blue", Stats: laststats.PlayerStats{Frags: 4}},
			},
		}

		expectTimestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", "2022-10-31 19:59:46 +0100")
		expect := lastscores.Entry{
			Timestamp: expectTimestamp,
			Mode:      "ffa",
			Participants: []lastscores.Participant{
				{Name: "XantoM", Frags: 6},
				{Name: "Xterm", Frags: 4},
			},
			Map: "dm6",
		}

		assert.Equal(t, expect, lastscores.FromLastStatsEntry(laststatsEntry))
	})
}
