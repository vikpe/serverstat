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
		players := []qclient.Client{}
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

	t.Run("2on2", func(t *testing.T) {
		players := []qclient.Client{
			{Name: qstring.New("XantoM"), Team: qstring.New("red")},
			{Name: qstring.New("Xterm"), Team: qstring.New("blue")},
			{Name: qstring.New("valla"), Team: qstring.New("blue")},
		}
		settings := qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "teamplay": "2", "map": "dm6"}
		assert.Equal(t, "2on2: blue (valla, Xterm) vs red (XantoM) [dm6]", qtitle.New(settings, players))
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

func TestParseMatchtag(t *testing.T) {
	testCases := map[string]string{
		"kombat":        "kombat",
		"testing stuff": "",
		"pausable game": "",
		"pause":         "",
		"":              "",
		"xy":            "",
		"xyz":           "xyz",
		"2on2":          "",
		"kombat 2on2":   "kombat 2on2",
	}

	for matchtag, expect := range testCases {
		t.Run(matchtag, func(t *testing.T) {
			assert.Equal(t, expect, qtitle.ParseMatchtag(matchtag))
		})
	}
}
