package analyze

import (
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qclient"
	"golang.org/x/exp/slices"
)

func HasSpectator(server mvdsv.Mvdsv, name string) bool {
	return HasQtvSpectator(server, name) || HasServerSpectator(server, name)
}

func HasQtvSpectator(server mvdsv.Mvdsv, name string) bool {
	return ListOfNamesContainsName(server.QtvStream.SpectatorNames, name)
}

func HasServerSpectator(server mvdsv.Mvdsv, name string) bool {
	return ListOfNamesContainsName(server.SpectatorNames, name)
}

func HasPlayer(server mvdsv.Mvdsv, name string) bool {
	return ListOfNamesContainsName(GetPlayerNames(server), name)
}

func GetPlayerNames(server mvdsv.Mvdsv) []string {
	playerNames := make([]string, 0)

	for _, player := range server.Players {
		playerNames = append(playerNames, player.Name.ToPlainString())
	}

	return playerNames
}

func HasClient(server mvdsv.Mvdsv, name string) bool {
	return HasPlayer(server, name) || HasSpectator(server, name)
}

func ListOfNamesContainsName(playerNames []string, name string) bool {
	if 0 == len(playerNames) {
		return false
	}

	normalizeNames := make([]string, 0)

	for _, playerName := range playerNames {
		normalizeNames = append(normalizeNames, normalizeName(playerName))
	}

	return slices.Contains(normalizeNames, name)
}

func normalizeName(name string) string {
	return strings.ToLower(name)
}

func IsIdle(server mvdsv.Mvdsv) bool {
	if 0 == server.PlayerSlots.Used {
		return true
	}

	if !server.Status.IsStandby() || server.Mode.IsRace() {
		return false
	}

	if server.Mode.IsXonX() && server.PlayerSlots.Free > 0 {
		return false
	}

	minIdleLimit := 3
	maxIdleLimit := 10
	idleLimit := clampInt(int(float64(server.PlayerSlots.Used)*1.5), minIdleLimit, maxIdleLimit)

	return MinPlayerTime(server.Players) >= idleLimit
}

func MinPlayerTime(players []qclient.Client) int {
	playerCount := len(players)

	if 0 == playerCount {
		return 0
	} else if 1 == playerCount {
		return players[0].Time
	}

	result := players[0].Time

	for _, p := range players {
		if p.Time < result {
			result = p.Time
		}
	}

	return result
}

func clampInt(value int, min int, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
