package qmode_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

func TestParse(t *testing.T) {
	type testCase struct {
		mode     string
		settings qsettings.Settings
	}

	testCases := []testCase{
		{"fortress", qsettings.Settings{"*gamedir": "fortress"}},
		{"race", qsettings.Settings{"*gamedir": "qw", "ktxmode": "race"}},
		{"coop", qsettings.Settings{"*gamedir": "qw", "teamplay": "2", "maxclients": "26"}},
		{"coop", qsettings.Settings{"*gamedir": "qw", "teamplay": "2", "maxclients": "12"}},
		{"1on1", qsettings.Settings{"*gamedir": "qw", "maxclients": "2"}},
		{"2on2", qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "teamplay": "2"}},
		{"4on4", qsettings.Settings{"*gamedir": "qw", "maxclients": "8", "teamplay": "2"}},
		{"ctf", qsettings.Settings{"*gamedir": "qw", "maxclients": "16", "teamplay": "4"}},
		{"ffa", qsettings.Settings{"*gamedir": "qw", "maxclients": "8"}},
		{"ffa", qsettings.Settings{"*gamedir": "ktx-ffa"}},
	}

	for _, tc := range testCases {
		t.Run(tc.mode, func(t *testing.T) {
			expect := qmode.Mode(tc.mode)
			assert.Equal(t, expect, qmode.Parse(tc.settings), tc.mode)
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
