package qmode

import (
	"fmt"
	"strings"

	"github.com/vikpe/serverstat/qserver/qsettings"
)

type Mode string

func (m Mode) Is(name string) bool {
	return strings.ToLower(name) == strings.ToLower(string(m))
}
func (m Mode) Is1on1() bool     { return m.Is("1on1") }
func (m Mode) Is2on2() bool     { return m.Is("2on2") }
func (m Mode) Is3on3() bool     { return m.Is("3on3") }
func (m Mode) Is4on4() bool     { return m.Is("4on4") }
func (m Mode) IsRace() bool     { return m.Is("race") }
func (m Mode) IsFfa() bool      { return m.Is("ffa") }
func (m Mode) IsCtf() bool      { return m.Is("ctf") }
func (m Mode) IsCoop() bool     { return m.Is("coop") }
func (m Mode) IsCustom() bool   { return m.Is("custom") }
func (m Mode) IsFortress() bool { return m.Is("fortress") }
func (m Mode) IsUnknown() bool  { return m.Is("unknown") }

func Parse(settings qsettings.Settings) Mode {
	gameDir := strings.ToLower(settings.Get("*gamedir", "unknown"))

	if "qw" != gameDir {
		return Mode(gameDir)
	}

	if settings.Has("ktxmode") {
		return Mode(strings.ToLower(settings.Get("ktxmode", "unknown")))
	}

	teamplay := settings.GetInt("teamplay", 0)
	maxClients := settings.GetInt("maxclients", 0)

	if teamplay > 0 {
		if 2 == teamplay && 26 == maxClients {
			return "coop"
		} else {
			playersPerTeam := maxClients / 2
			modeName := fmt.Sprintf("%don%d", playersPerTeam, playersPerTeam)
			return Mode(modeName)
		}
	}

	if 2 == maxClients {
		return "1on1"
	} else {
		return "ffa"
	}
}
