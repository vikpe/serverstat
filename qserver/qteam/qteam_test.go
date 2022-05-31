package qteam_test

import (
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
	xantom := qclient.Client{
		Name:   qstring.New("XantoM"),
		Team:   qstring.New("f0m"),
		Skin:   "",
		Colors: [2]uint8{4, 2},
		Frags:  0,
		Ping:   0,
		Time:   0,
		CC:     "",
	}
	bps := qclient.Client{
		Name:   qstring.New("bps"),
		Team:   qstring.New("-s-"),
		Skin:   "",
		Colors: [2]uint8{4, 2},
		Frags:  0,
		Ping:   0,
		Time:   0,
		CC:     "",
	}
	valla := qclient.Client{
		Name:   qstring.New("valla"),
		Team:   qstring.New("f0m"),
		Skin:   "",
		Colors: [2]uint8{4, 2},
		Frags:  0,
		Ping:   0,
		Time:   0,
		CC:     "",
	}
	players := []qclient.Client{xantom, bps, valla}
	teams := qteam.FromPlayers(players)

	expect := []qteam.Team{
		{Name: qstring.New("-s-"), Players: []qclient.Client{bps}},
		{Name: qstring.New("f0m"), Players: []qclient.Client{xantom, valla}},
	}
	assert.Equal(t, expect, teams)
}

func TestExport(t *testing.T) {
	player1 := qclient.Client{Name: qstring.New("XantoM"), Colors: [2]uint8{4, 2}, Frags: 12}
	player2 := qclient.Client{Name: qstring.New("Milton"), Colors: [2]uint8{4, 2}, Frags: 8}
	player3 := qclient.Client{Name: qstring.New("bps"), Colors: [2]uint8{13, 5}, Frags: 8}

	team := qteam.Team{
		Name:    qstring.New("red"),
		Players: []qclient.Client{player1, player2, player3},
	}
	expect := qteam.TeamExport{
		Name:      qstring.New("red"),
		NameColor: "www",
		Frags:     28,
		Colors:    [2]uint8{4, 2},
		Players:   []qclient.Client{player1, player3, player2},
	}
	assert.Equal(t, expect, qteam.Export(team))
}
