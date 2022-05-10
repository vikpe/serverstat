package qmode_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

func TestParse(t *testing.T) {
	testCases := map[string]qsettings.Settings{
		"fortress": {"*gamedir": "fortress"},
		"race":     {"*gamedir": "qw", "ktxmode": "race"},
		"coop":     {"*gamedir": "qw", "teamplay": "2", "maxclients": "26"},
		"1on1":     {"*gamedir": "qw", "maxclients": "2"},
		"2on2":     {"*gamedir": "qw", "maxclients": "4", "teamplay": "2"},
		"4on4":     {"*gamedir": "qw", "maxclients": "8", "teamplay": "2"},
		"ffa":      {"*gamedir": "qw", "maxclients": "8"},
	}

	for expectedModeName, settings := range testCases {
		t.Run(expectedModeName, func(t *testing.T) {
			expect := qmode.Mode(expectedModeName)
			assert.Equal(t, expect, qmode.Parse(settings), expectedModeName)
		})
	}
}

func TestModeValidators(t *testing.T) {
	testCases := map[string]func(m qmode.Mode) bool{
		"1on1":     func(m qmode.Mode) bool { return m.Is1on1() },
		"2on2":     func(m qmode.Mode) bool { return m.Is2on2() },
		"3on3":     func(m qmode.Mode) bool { return m.Is3on3() },
		"4on4":     func(m qmode.Mode) bool { return m.Is4on4() },
		"race":     func(m qmode.Mode) bool { return m.IsRace() },
		"ffa":      func(m qmode.Mode) bool { return m.IsFfa() },
		"ctf":      func(m qmode.Mode) bool { return m.IsCtf() },
		"coop":     func(m qmode.Mode) bool { return m.IsCoop() },
		"fortress": func(m qmode.Mode) bool { return m.IsFortress() },
		"unknown":  func(m qmode.Mode) bool { return m.IsUnknown() },
	}

	for currentModeName, currentValidator := range testCases {
		for modeName, otherValidator := range testCases {
			mode := qmode.Mode(currentModeName)

			if modeName == currentModeName {
				assert.True(t, currentValidator(mode))
				assert.True(t, currentValidator(qmode.Mode(strings.ToUpper(currentModeName))))
			} else {
				assert.False(t, otherValidator(mode))
			}
		}
	}
}
