package qtitle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtitle"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestMvdsv_Title(t *testing.T) {
	t.Run("matchtag", func(t *testing.T) {
		players := []qclient.Client{}
		settings := qsettings.Settings{"matchtag": "kombat"}
		assert.Equal(t, "kombat / unknown []", qtitle.New(settings, players))
	})

	t.Run("ffa", func(t *testing.T) {
		players := []qclient.Client{{Name: qstring.New("XantoM")}}
		settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "map": "dm2"}
		assert.Equal(t, "ffa [dm2]", qtitle.New(settings, players))
	})

	t.Run("1on1", func(t *testing.T) {
		players := []qclient.Client{
			{Name: qstring.New("XantoM")},
			{Name: qstring.New("Xterm")},
		}
		settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "2", "map": "dm6"}
		assert.Equal(t, "1on1: XantoM vs Xterm [dm6]", qtitle.New(settings, players))
	})

	t.Run("xonx", func(t *testing.T) {
		t.Run("no players", func(t *testing.T) {
			players := []qclient.Client{}
			settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "teamplay": "2", "map": "dm6"}
			assert.Equal(t, "2on2 [dm6]", qtitle.New(settings, players))
		})

		t.Run("one team/player", func(t *testing.T) {
			players := []qclient.Client{
				{Name: qstring.New("XantoM"), Team: qstring.New("red")},
			}
			settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "teamplay": "2", "map": "dm6"}
			assert.Equal(t, "2on2: red (XantoM) [dm6]", qtitle.New(settings, players))
		})

		t.Run("<= 2 teams", func(t *testing.T) {
			players := []qclient.Client{
				{Name: qstring.New("XantoM"), Team: qstring.New("red")},
				{Name: qstring.New("Xterm"), Team: qstring.New("blue")},
				{Name: qstring.New("valla"), Team: qstring.New("blue")},
			}
			settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "teamplay": "2", "map": "dm6"}
			assert.Equal(t, "2on2: blue (valla, Xterm) vs red (XantoM) [dm6]", qtitle.New(settings, players))
		})

		t.Run(">2 teams", func(t *testing.T) {
			players := []qclient.Client{
				{Name: qstring.New("XantoM"), Team: qstring.New("red")},
				{Name: qstring.New("Xterm"), Team: qstring.New("blue")},
				{Name: qstring.New("valla"), Team: qstring.New("mix")},
			}
			settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "teamplay": "2", "map": "dm6"}
			assert.Equal(t, "2on2: valla, XantoM, Xterm [dm6]", qtitle.New(settings, players))
		})

		t.Run("many teams", func(t *testing.T) {
			players := []qclient.Client{
				{Name: qstring.New("hangtime"), Team: qstring.New("+er+")},
				{Name: qstring.New("FU-hto"), Team: qstring.New("-fu-")},
				{Name: qstring.New("alice"), Team: qstring.New("1")},
				{Name: qstring.New("NinJaA"), Team: qstring.New("blue")},
				{Name: qstring.New("sniegov"), Team: qstring.New("blue")},
				{Name: qstring.New("Xterm"), Team: qstring.New("com")},
				{Name: qstring.New("eclip"), Team: qstring.New("r0t")},
			}
			settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "8", "teamplay": "2", "map": "dm3"}
			assert.Equal(t, "4on4: alice, eclip, FU-hto, hangtime, NinJaA, sniegov, Xterm [dm3]", qtitle.New(settings, players))
		})
	})

	t.Run("coop", func(t *testing.T) {
		players := []qclient.Client{
			{Name: qstring.New("Xterm")},
			{Name: qstring.New("andeh")},
			{Name: qstring.New("XantoM")},
		}
		settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "26", "teamplay": "2", "map": "bloodfest"}
		assert.Equal(t, "coop: andeh, XantoM, Xterm [bloodfest]", qtitle.New(settings, players))
	})
}

func BenchmarkNew(b *testing.B) {
	players := []qclient.Client{
		{Name: qstring.New("hangtime"), Team: qstring.New("+er+")},
		{Name: qstring.New("FU-hto"), Team: qstring.New("-fu-")},
		{Name: qstring.New("alice"), Team: qstring.New("1")},
		{Name: qstring.New("NinJaA"), Team: qstring.New("blue")},
		{Name: qstring.New("sniegov"), Team: qstring.New("blue")},
		{Name: qstring.New("Xterm"), Team: qstring.New("com")},
		{Name: qstring.New("eclip"), Team: qstring.New("r0t")},
	}
	settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "8", "teamplay": "2", "map": "dm3", "matchtag": "kombat"}

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("teamplay", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qtitle.New(settings, players)
		}
	})

	delete(settings, "teamplay")

	b.Run("no teamplay", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qtitle.New(settings, players)
		}
	})

	b.Run("no players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qtitle.New(settings, make([]qclient.Client, 0))
		}
	})
}

func BenchmarkTeamCount(b *testing.B) {
	players := []qclient.Client{
		{Name: qstring.New("hangtime"), Team: qstring.New("+er+")},
		{Name: qstring.New("FU-hto"), Team: qstring.New("-fu-")},
		{Name: qstring.New("alice"), Team: qstring.New("1")},
		{Name: qstring.New("NinJaA"), Team: qstring.New("blue")},
		{Name: qstring.New("sniegov"), Team: qstring.New("blue")},
		{Name: qstring.New("Xterm"), Team: qstring.New("com")},
		{Name: qstring.New("eclip"), Team: qstring.New("r0t")},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qtitle.TeamCount(players)
	}
}

func TestParseMatchtag(t *testing.T) {
	testCases := map[string]string{
		// too short
		"":     "",
		"xy":   "",
		" xy ": "",

		// ignored phrases
		"testing stuff": "",
		"pausable game": "",
		"pause":         "",
		"2on2":          "",
		"race":          "",

		// symbols
		"--great!!--":        "great",
		"--great-stuff--":    "great-stuff",
		"***":                "",
		"..ab":               "",
		"acb._[]{}()*\\/def": "acb def",
		"..abc":              "abc",
		".a.b.":              "a b",

		// mixed cases
		"[2on2]":  "",
		"race...": "",

		// keep
		"xyz":         "xyz",
		"kombat":      "kombat",
		"kombat 2on2": "kombat 2on2",
	}

	for matchtag, expect := range testCases {
		t.Run(matchtag, func(t *testing.T) {
			assert.Equal(t, expect, qtitle.ParseMatchtag(matchtag))
		})
	}
}

func BenchmarkParseMatchtag(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("long value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qtitle.ParseMatchtag("acb._[]{}()*\\/def")
		}
	})

	b.Run("short value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qtitle.ParseMatchtag("..abc")
		}
	})
}

func TestTeamCount(t *testing.T) {
	t.Run("no players", func(t *testing.T) {
		assert.Equal(t, 0, qtitle.TeamCount(nil))
	})

	t.Run("one player", func(t *testing.T) {
		players := []qclient.Client{{Team: qstring.New("1")}}
		assert.Equal(t, 1, qtitle.TeamCount(players))
	})

	t.Run("many players", func(t *testing.T) {
		players := []qclient.Client{
			{Team: qstring.New("1")},
			{Team: qstring.New("2")},
			{Team: qstring.New("3")},
			{Team: qstring.New("4")},
			{Team: qstring.New("5")},
			{Team: qstring.New("6")},
		}
		assert.Equal(t, 6, qtitle.TeamCount(players))
	})

	t.Run("many players - one team", func(t *testing.T) {
		players := []qclient.Client{
			{Team: qstring.New("1")},
			{Team: qstring.New("1")},
			{Team: qstring.New("1")},
		}
		assert.Equal(t, 1, qtitle.TeamCount(players))
	})
}
