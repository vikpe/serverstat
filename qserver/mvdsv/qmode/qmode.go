package qmode

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

const (
	mode10on10    = "10on10"
	mode1on1      = "1on1"
	mode2on2      = "2on2"
	mode3on3      = "3on3"
	mode4on4      = "4on4"
	modeClanArena = "clan arena"
	modeCoop      = "coop"
	modeCtf       = "ctf"
	modeFfa       = "ffa"
	modeFortress  = "fortress"
	modeRace      = "race"
	modeUnknown   = "unknown"
)

type Mode string

func (m Mode) Is(name string) bool {
	return strings.ToLower(name) == strings.ToLower(string(m))
}
func (m Mode) Is10on10() bool    { return m.Is(mode10on10) }
func (m Mode) Is1on1() bool      { return m.Is(mode1on1) }
func (m Mode) Is2on2() bool      { return m.Is(mode2on2) }
func (m Mode) Is3on3() bool      { return m.Is(mode3on3) }
func (m Mode) Is4on4() bool      { return m.Is(mode4on4) }
func (m Mode) IsClanArena() bool { return m.Is(modeClanArena) }
func (m Mode) IsCoop() bool      { return m.Is(modeCoop) }
func (m Mode) IsCtf() bool       { return m.Is(modeCtf) }
func (m Mode) IsFfa() bool       { return m.Is(modeFfa) }
func (m Mode) IsFortress() bool  { return m.Is(modeFortress) }
func (m Mode) IsRace() bool      { return m.Is(modeRace) }
func (m Mode) IsUnknown() bool   { return m.Is(modeUnknown) }
func (m Mode) IsCustom() bool    { return !(m.IsFfa() || m.IsXonX()) }
func (m Mode) IsXonX() bool {
	xonxModes := []string{mode1on1, mode2on2, mode3on3, mode4on4, mode10on10}
	return lo.Contains(xonxModes, string(m))
}

func Parse(settings qsettings.Settings) (Mode, string) {
	gameDir := strings.ToLower(settings.Get("*gamedir", "qw"))

	// check gamedir
	customGameDirs := map[string]string{
		"ktx-ffa": modeFfa,
	}

	if modeName, ok := customGameDirs[gameDir]; ok {
		return Mode(modeName), ""
	}

	if "qw" != gameDir {
		return Mode(gameDir), ""
	}

	// check mode and ktx mode
	if settingsMode, ok := settings["mode"]; ok {
		if strings.HasSuffix(settingsMode, modeRace) {
			return modeRace, ""
		}

		settingsMode := strings.ToLower(settingsMode)

		if strings.Contains(settingsMode, "-") {
			parts := strings.SplitN(settingsMode, "-", 2)
			return Mode(parts[0]), parts[1]
		}

		return Mode(settingsMode), ""
	}

	if ktxMode, ok := settings["ktxmode"]; ok {
		return Mode(strings.ToLower(ktxMode)), ""
	}

	// derive from settings
	teamplay := settings.GetInt("teamplay", 0)
	maxClients := settings.GetInt("maxclients", 0)

	if teamplay > 0 {
		deathmatch := settings.GetInt("deathmatch", 0)

		if 2 == teamplay && lo.Contains([]int{26, 24, 12}, maxClients) {
			return modeCoop, ""
		} else if 4 == teamplay {
			if 16 == maxClients {
				return modeCtf, ""
			}

			if 5 == deathmatch {
				return modeClanArena, ""
			}
		}

		playersPerTeam := maxClients / 2
		modeName := fmt.Sprintf("%don%d", playersPerTeam, playersPerTeam)
		return Mode(modeName), ""
	}

	if 2 == maxClients {
		return mode1on1, ""
	}

	return modeFfa, ""
}
