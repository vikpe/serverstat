package qtitle

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qteam"
)

func New(settings qsettings.Settings, players []qclient.Client) string {
	title := ""

	// matchtag
	matchTag := ParseMatchtag(settings.Get("matchtag", ""))

	if matchTag != "" {
		title += fmt.Sprintf("%s / ", matchTag)
	}

	// mode
	mode := qmode.Parse(settings)
	title += string(mode)

	// participants
	participants := make([]string, 0)
	isTeamplay := settings.GetInt("teamplay", 0) > 0

	if isTeamplay && !mode.IsCoop() {
		teams := qteam.FromPlayers(players)

		for _, t := range teams {
			participants = append(participants, t.String())
		}
	} else if !mode.IsFfa() {
		for _, p := range players {
			participants = append(participants, p.Name.ToPlainString())
		}
	}

	if len(participants) > 0 {
		var participantDelimiter string

		if mode.IsCoop() {
			participantDelimiter = ", "
		} else {
			participantDelimiter = " vs "
		}

		sort.Slice(participants, func(i, j int) bool {
			return strings.ToLower(participants[i]) < strings.ToLower(participants[j])
		})
		title += ": " + strings.Join(participants, participantDelimiter)
	}

	// map
	title += fmt.Sprintf(" [%s]", settings.Get("map", ""))

	return title
}

func ParseMatchtag(matchtag string) string {
	const minLength = 3

	if len(matchtag) < minLength {
		return ""
	}

	ignoreList := []string{
		"pause",
		"pausable",
		"test",
	}

	for _, word := range ignoreList {
		if strings.Contains(matchtag, word) {
			return ""
		}
	}

	return matchtag
}
