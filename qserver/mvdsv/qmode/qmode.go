package qmode

import (
	"fmt"
	"strings"

	"github.com/vikpe/serverstat/qserver/qsettings"
	"golang.org/x/exp/slices"
)

const (
	mode1on1     = "1on1"
	mode2on2     = "2on2"
	mode3on3     = "3on3"
	mode4on4     = "4on4"
	modeRace     = "Race"
	modeFfa      = "ffa"
	modeCtf      = "ctf"
	modeCoop     = "coop"
	modeFortress = "fortress"
	modeUnknown  = "unknown"
)

type Mode string

func (m Mode) Is(name string) bool {
	return strings.ToLower(name) == strings.ToLower(string(m))
}
func (m Mode) Is1on1() bool     { return m.Is(mode1on1) }
func (m Mode) Is2on2() bool     { return m.Is(mode2on2) }
func (m Mode) Is3on3() bool     { return m.Is(mode3on3) }
func (m Mode) Is4on4() bool     { return m.Is(mode4on4) }
func (m Mode) IsRace() bool     { return m.Is(modeRace) }
func (m Mode) IsFfa() bool      { return m.Is(modeFfa) }
func (m Mode) IsCtf() bool      { return m.Is(modeCtf) }
func (m Mode) IsCoop() bool     { return m.Is(modeCoop) }
func (m Mode) IsFortress() bool { return m.Is(modeFortress) }
func (m Mode) IsUnknown() bool  { return m.Is(modeUnknown) }

func Parse(settings qsettings.Settings) Mode {
	gameDir := strings.ToLower(settings.Get("*gamedir", modeUnknown))

	// check gamedir
	customGameDirs := map[string]string{
		"ktx-ffa": modeFfa,
	}

	if modeName, ok := customGameDirs[gameDir]; ok {
		return Mode(modeName)
	}

	// check ktx mode
	if "qw" != gameDir {
		return Mode(gameDir)
	}

	if ktxMode, ok := settings["ktxmode"]; ok {
		return Mode(strings.ToLower(ktxMode))
	}

	// derive from settings
	teamplay := settings.GetInt("teamplay", 0)
	maxClients := settings.GetInt("maxclients", 0)

	if teamplay > 0 {
		if 2 == teamplay && slices.Contains([]int{26, 24, 12}, maxClients) {
			return modeCoop
		} else if 4 == teamplay && 16 == maxClients {
			return modeCtf
		} else {
			playersPerTeam := maxClients / 2
			modeName := fmt.Sprintf("%don%d", playersPerTeam, playersPerTeam)
			return Mode(modeName)
		}
	}

	if 2 == maxClients {
		return mode1on1
	} else {
		return modeFfa
	}
}
