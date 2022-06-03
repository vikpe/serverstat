package qtitle

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qutil"
)

func New(settings qsettings.Settings, players []qclient.Client) string {
	title := strings.Builder{}

	// matchtag
	matchTag := ParseMatchtag(settings.Get("matchtag", ""))

	if matchTag != "" {
		title.WriteString(fmt.Sprintf("%s / ", matchTag))
	}

	// mode
	mode := qmode.Parse(settings)
	title.WriteString(string(mode))

	if 0 == len(players) || mode.IsFfa() {
		title.WriteString(fmtMap(settings.Get("map", "")))
		return title.String()
	}

	// participants
	participants := make([]string, 0)
	var participantDelimiter string
	isTeamplay := settings.GetInt("teamplay", 0) > 0

	if isTeamplay && !mode.IsCoop() && TeamCount(players) <= 2 {
		teams := qteam.FromPlayers(players)

		participantDelimiter = " vs "

		for _, t := range teams {
			participants = append(participants, t.String())
		}
	} else if !mode.IsFfa() {
		for _, p := range players {
			participants = append(participants, p.Name.ToPlainString())
		}

		if mode.Is1on1() {
			participantDelimiter = " vs "
		} else {
			participantDelimiter = ", "
		}
	}

	sort.Slice(participants, func(i, j int) bool {
		return strings.ToLower(participants[i]) < strings.ToLower(participants[j])
	})
	title.WriteString(": " + strings.Join(participants, participantDelimiter))

	// map
	title.WriteString(fmtMap(settings.Get("map", "")))

	return title.String()
}

func fmtMap(value string) string {
	return fmt.Sprintf(" [%s]", value)
}

func TeamCount(players []qclient.Client) int {
	if len(players) < 2 {
		return len(players)
	}

	teams := make(map[string]bool, 0)

	for _, p := range players {
		teams[p.Team.ToPlainString()] = true
	}

	return len(teams)
}

func ParseMatchtag(matchtag string) string {
	result := qutil.TrimSymbols(matchtag)

	if len(result) < 3 {
		return ""
	}

	ignorePartial := []string{
		"pause",
		"pausable",
		"test",
	}

	for _, needle := range ignorePartial {
		if strings.Contains(result, needle) {
			return ""
		}
	}

	ignoreEqual := []string{
		"duel",
		"1on1",
		"2on2",
		"3on3",
		"4on4",
		"ffa",
		"ctf",
		"race",
	}

	for _, word := range ignoreEqual {
		if result == word {
			return ""
		}
	}

	return result
}
