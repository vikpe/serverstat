package bot

import (
	"strings"
)

func IsBotName(name string) bool {
	if 0 == len(name) {
		return false
	}

	knownBotNames := []string{
		"[ServeMe]",
		"twitch.tv/vikpe",
	}

	return strings.Contains(strings.Join(knownBotNames, "\""), name)
}

func IsBotPing(ping int) bool {
	return 10 == ping || ping < -400
}
