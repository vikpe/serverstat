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
		{"1on1", qsettings.Settings{"*gamedir": "qw", "maxclients": "2"}},
		{"2on2", qsettings.Settings{"*gamedir": "qw", "maxclients": "4", "teamplay": "2"}},
		{"4on4", qsettings.Settings{"*gamedir": "qw", "maxclients": "8", "teamplay": "2"}},
		{"clan arena", qsettings.Settings{"*gamedir": "qw", "deathmatch": "5", "teamplay": "4"}},
		{"coop", qsettings.Settings{"*gamedir": "qw", "teamplay": "2", "maxclients": "12"}},
		{"coop", qsettings.Settings{"*gamedir": "qw", "teamplay": "2", "maxclients": "24"}},
		{"coop", qsettings.Settings{"*gamedir": "qw", "teamplay": "2", "maxclients": "26"}},
		{"ctf", qsettings.Settings{"*gamedir": "qw", "maxclients": "16", "teamplay": "4"}},
		{"ffa", qsettings.Settings{"*gamedir": "ktx-ffa"}},
		{"ffa", qsettings.Settings{"*gamedir": "qw", "maxclients": "8"}},
		{"fortress", qsettings.Settings{"*gamedir": "fortress"}},
		{"race", qsettings.Settings{"*gamedir": "qw", "ktxmode": "race"}},
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
		"10on10":     func(m qmode.Mode) bool { return m.Is10on10() },
		"1on1":       func(m qmode.Mode) bool { return m.Is1on1() },
		"2on2":       func(m qmode.Mode) bool { return m.Is2on2() },
		"3on3":       func(m qmode.Mode) bool { return m.Is3on3() },
		"4on4":       func(m qmode.Mode) bool { return m.Is4on4() },
		"clan arena": func(m qmode.Mode) bool { return m.IsClanArena() },
		"coop":       func(m qmode.Mode) bool { return m.IsCoop() },
		"ctf":        func(m qmode.Mode) bool { return m.IsCtf() },
		"ffa":        func(m qmode.Mode) bool { return m.IsFfa() },
		"fortress":   func(m qmode.Mode) bool { return m.IsFortress() },
		"race":       func(m qmode.Mode) bool { return m.IsRace() },
		"unknown":    func(m qmode.Mode) bool { return m.IsUnknown() },
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

	assert.True(t, qmode.Mode("1on1").IsXonX())
	assert.True(t, qmode.Mode("4on4").IsXonX())
	assert.False(t, qmode.Mode("ffa").IsXonX())
	assert.False(t, qmode.Mode("coop").IsXonX())
	assert.False(t, qmode.Mode("race").IsXonX())

	assert.False(t, qmode.Mode("1on1").IsCustom())
	assert.False(t, qmode.Mode("4on4").IsCustom())
	assert.False(t, qmode.Mode("ffa").IsCustom())
	assert.True(t, qmode.Mode("coop").IsCustom())
	assert.True(t, qmode.Mode("race").IsCustom())
}
