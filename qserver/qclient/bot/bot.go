package bot

import (
	"github.com/samber/lo"
)

func IsBotName(name string) bool {
	if 0 == len(name) {
		return false
	}

	knownBotNames := []string{
		"[ServeMe]",
		"twitch.tv/vikpe",
	}

	return lo.Contains(knownBotNames, name)
}

func IsBotPing(ping int) bool {
	return 10 == ping
}
