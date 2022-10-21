package qteam_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestTeam_Frags(t *testing.T) {
	team := qteam.Team{
		Name: qstring.New("red"),
		Players: []qclient.Client{
			{Frags: 5},
			{Frags: -2},
			{Frags: 0},
		},
	}
	assert.Equal(t, 3, team.Frags())
}

func TestTeam_Ping(t *testing.T) {
	t.Run("no players", func(t *testing.T) {
		team := qteam.Team{
			Name:    qstring.New("red"),
			Players: []qclient.Client{},
		}
		assert.Equal(t, 0, team.Ping())
	})

	t.Run("single player", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("red"),
			Players: []qclient.Client{
				{Ping: 5},
			},
		}
		assert.Equal(t, 5, team.Ping())
	})

	t.Run("multiple players/bots", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("red"),
			Players: []qclient.Client{
				{Ping: 12},
				{Ping: 10}, // bot
				{Ping: 24},
			},
		}
		assert.Equal(t, 18, team.Ping())
	})
}

func TestTeam_Colors(t *testing.T) {
	t.Run("no players", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("red"),
		}
		assert.Equal(t, [2]uint8{0, 0}, team.Colors())
	})

	t.Run("majority colors", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("red"),
			Players: []qclient.Client{
				{Colors: [2]uint8{0, 0}},
				{Colors: [2]uint8{4, 2}},
				{Colors: [2]uint8{4, 2}},
			},
		}
		assert.Equal(t, [2]uint8{4, 2}, team.Colors())
	})

	t.Run("no majority colors (use lowst color) [1]", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("red"),
			Players: []qclient.Client{
				{Colors: [2]uint8{9, 6}},
				{Colors: [2]uint8{2, 0}},
				{Colors: [2]uint8{1, 11}},
				{Colors: [2]uint8{4, 2}},
			},
		}
		assert.Equal(t, [2]uint8{1, 11}, team.Colors())
	})

	t.Run("no majority colors (use lowest color) [2]", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("red"),
			Players: []qclient.Client{
				{Colors: [2]uint8{4, 2}},
				{Colors: [2]uint8{0, 5}},
				{Colors: [2]uint8{0, 5}},
				{Colors: [2]uint8{4, 2}},
			},
		}
		assert.Equal(t, [2]uint8{0, 5}, team.Colors())
	})
}

func TestTeam_String(t *testing.T) {
	t.Run("empty team", func(t *testing.T) {
		team := qteam.Team{Name: qstring.New("mix")}
		assert.Equal(t, "mix", team.String())
	})

	t.Run("has few players", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("f0m"),
			Players: []qclient.Client{
				{Name: qstring.New("XantoM")},
				{Name: qstring.New("valla")},
			},
		}
		assert.Equal(t, "f0m (valla, XantoM)", team.String())
	})

	t.Run("has a lot of players", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("f0m"),
			Players: []qclient.Client{
				{Name: qstring.New("XantoM 1")},
				{Name: qstring.New("XantoM 2")},
				{Name: qstring.New("XantoM 3")},
				{Name: qstring.New("XantoM 4")},
				{Name: qstring.New("XantoM 5")},
			},
		}
		assert.Equal(t, "f0m", team.String())
	})

	t.Run("don't strip prefix/suffix for single player", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("oeks"),
			Players: []qclient.Client{
				{Name: qstring.New("nig......axe")},
			},
		}
		assert.Equal(t, "oeks (nig......axe)", team.String())
	})

	t.Run("strip common prefix/suffix for multiple players", func(t *testing.T) {
		team := qteam.Team{
			Name: qstring.New("oeks"),
			Players: []qclient.Client{
				{Name: qstring.New("--nig......axe")},
				{Name: qstring.New("--trl......axe")},
			},
		}
		assert.Equal(t, "oeks (nig, trl)", team.String())
	})
}

func TestFromPlayers(t *testing.T) {
	teamless := qclient.Client{
		Name:   qstring.New("XantoM"),
		Team:   qstring.New(""),
		Colors: [2]uint8{0, 0},
		Frags:  10,
	}
	xantom := qclient.Client{
		Name:   qstring.New("XantoM"),
		Team:   qstring.New("f0m"),
		Colors: [2]uint8{4, 2},
		Frags:  9,
	}
	xterm := qclient.Client{
		Name:   qstring.New("Xterm"),
		Team:   qstring.New("f0m"),
		Colors: [2]uint8{4, 2},
		Frags:  12,
	}
	bps := qclient.Client{
		Name:   qstring.New("bps"),
		Team:   qstring.New("-s-"),
		Colors: [2]uint8{4, 2},
		Frags:  3,
	}
	valla := qclient.Client{
		Name:   qstring.New("valla"),
		Team:   qstring.New("f0m"),
		Colors: [2]uint8{4, 2},
		Frags:  9,
	}
	players := []qclient.Client{xantom, xterm, bps, valla, teamless}
	expect := []qteam.Team{
		{Name: qstring.New(""), Players: []qclient.Client{teamless}},
		{Name: qstring.New("-s-"), Players: []qclient.Client{bps}},
		{Name: qstring.New("f0m"), Players: []qclient.Client{xterm, valla, xantom}},
	}

	assert.Equal(t, expect, qteam.FromPlayers(players))
	assert.Equal(t, []qclient.Client{xantom, xterm, bps, valla, teamless}, players) // ensure players slice is unchanged
}

func BenchmarkFromPlayers(b *testing.B) {
	xantom := qclient.Client{
		Name:   qstring.New("XantoM"),
		Team:   qstring.New("f0m"),
		Colors: [2]uint8{4, 2},
		Frags:  9,
	}
	xterm := qclient.Client{
		Name:   qstring.New("Xterm"),
		Team:   qstring.New("f0m"),
		Colors: [2]uint8{4, 2},
		Frags:  12,
	}
	bps := qclient.Client{
		Name:   qstring.New("bps"),
		Team:   qstring.New("-s-"),
		Colors: [2]uint8{4, 2},
		Frags:  3,
	}
	valla := qclient.Client{
		Name:   qstring.New("valla"),
		Team:   qstring.New("f0m"),
		Colors: [2]uint8{4, 2},
		Frags:  9,
	}
	paradoks := qclient.Client{
		Name:   qstring.New("paradoks"),
		Team:   qstring.New("]sr["),
		Colors: [2]uint8{3, 11},
		Frags:  23,
	}

	players := []qclient.Client{xantom, xterm, bps, valla, paradoks}

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("no players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qteam.FromPlayers(make([]qclient.Client, 0))
		}
	})

	b.Run("one players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qteam.FromPlayers(players[0:1])
		}
	})

	b.Run("many players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qteam.FromPlayers(players)
		}
	})
}

func TestExport(t *testing.T) {
	player1 := qclient.Client{Name: qstring.New("XantoM"), Colors: [2]uint8{4, 2}, Frags: 12, Ping: 12}
	player2 := qclient.Client{Name: qstring.New("Milton"), Colors: [2]uint8{4, 2}, Frags: 8, Ping: 14}
	player3 := qclient.Client{Name: qstring.New("bps"), Colors: [2]uint8{13, 5}, Frags: 8, Ping: 16}

	team := qteam.Team{
		Name:    qstring.New("red"),
		Players: []qclient.Client{player1, player2, player3},
	}
	expect := qteam.TeamExport{
		Name:      qstring.New("red"),
		NameColor: "www",
		Frags:     28,
		Ping:      14,
		Colors:    [2]uint8{4, 2},
		Players:   []qclient.Client{player1, player3, player2},
	}
	assert.Equal(t, expect, qteam.Export(team))
}

func TestTeam_MarshalJSON(t *testing.T) {
	player1 := qclient.Client{Name: qstring.New("XantoM"), Colors: [2]uint8{4, 2}, Frags: 12, Ping: 30}
	player2 := qclient.Client{Name: qstring.New("bps"), Colors: [2]uint8{13, 5}, Frags: 8, Ping: 20}
	team := qteam.Team{
		Name:    qstring.New("red"),
		Players: []qclient.Client{player1, player2},
	}
	jsonValue, _ := json.Marshal(team)
	expect := `{"name":"red","name_color":"www","frags":20,"ping":25,"colors":[4,2],"players":[{"name":"XantoM","name_color":"wwwwww","team":"","team_color":"","skin":"","colors":[4,2],"frags":12,"ping":30,"time":0,"cc":"","is_bot":false},{"name":"bps","name_color":"www","team":"","team_color":"","skin":"","colors":[13,5],"frags":8,"ping":20,"time":0,"cc":"","is_bot":false}]}`
	assert.Equal(t, expect, string(jsonValue))
}
