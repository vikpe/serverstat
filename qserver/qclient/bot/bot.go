package bot

import (
	"golang.org/x/exp/slices"
)

func IsBotName(name string) bool {
	if 0 == len(name) {
		return false
	}

	knownBotNames := []string{
		"[ServeMe]",
		"twitch.tv/vikpe",
	}

	return slices.Contains(knownBotNames, name)
}

func IsBotPing(ping int) bool {
	return 10 == ping || ping < -400 || ping > 400
}
