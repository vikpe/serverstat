package analyze

import (
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qclient"
)

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
	maxIdleLimit := 8
	idleLimit := clampInt(server.PlayerSlots.Used, minIdleLimit, maxIdleLimit)

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
