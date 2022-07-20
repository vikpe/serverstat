package analyze

import (
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/wildcard"
)

func HasSpectator(server mvdsv.Mvdsv, pattern string) bool {
	return HasQtvSpectator(server, pattern) || HasServerSpectator(server, pattern)
}

func HasQtvSpectator(server mvdsv.Mvdsv, pattern string) bool {
	return wildcard.MatchSliceCI(pattern, server.QtvStream.SpectatorNames)
}

func HasServerSpectator(server mvdsv.Mvdsv, pattern string) bool {
	return wildcard.MatchSliceCI(pattern, server.SpectatorNames)
}

func HasPlayer(server mvdsv.Mvdsv, pattern string) bool {
	return wildcard.MatchSliceCI(pattern, GetPlayerPlainNames(server))
}

func GetPlayerPlainNames(server mvdsv.Mvdsv) []string {
	playerNames := make([]string, 0)

	for _, player := range server.Players {
		playerNames = append(playerNames, player.Name.ToPlainString())
	}

	return playerNames
}

func HasClient(server mvdsv.Mvdsv, name string) bool {
	return HasPlayer(server, name) || HasSpectator(server, name)
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
	idleLimit := qutil.ClampInt(int(float64(server.PlayerSlots.Used)*1.5), minIdleLimit, maxIdleLimit)

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

func IsSpeccable(server mvdsv.Mvdsv) bool {
	if len(server.QtvStream.Url) > 0 {
		return true
	}

	return server.SpectatorSlots.Free > 0 && !RequiresPassword(server)
}

func RequiresPassword(server mvdsv.Mvdsv) bool {
	needpass := server.Settings.GetInt("needpass", 0)

	if 0 == needpass {
		return false
	}
	const spectatorPasswordBit = 2
	return (needpass & spectatorPasswordBit) > 0
}
