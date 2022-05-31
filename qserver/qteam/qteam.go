package qteam

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/serverstat/qutil"
)

type Team struct {
	Name    qstring.QuakeString
	Players []qclient.Client
}

type TeamExport struct {
	Name      qstring.QuakeString
	NameColor string
	Frags     int
	Colors    [2]uint8
	Players   []qclient.Client
}

func Export(t Team) TeamExport {
	qclient.SortPlayers(t.Players)

	return TeamExport{
		Name:      t.Name,
		NameColor: t.Name.ToColorCodes(),
		Colors:    t.Colors(),
		Frags:     t.Frags(),
		Players:   t.Players,
	}
}

func (t Team) MarshalJSON() ([]byte, error) {
	return qutil.MarshalNoEscapeHtml(Export(t))
}

func (t Team) Frags() int {
	frags := 0

	for _, p := range t.Players {
		frags += p.Frags
	}

	return frags
}

func (t Team) Colors() [2]uint8 {
	if 0 == len(t.Players) {
		return [2]uint8{0, 0}
	}

	colorCount := make(map[[2]uint8]int, 0)

	for _, p := range t.Players {
		colorCount[p.Colors]++
	}

	isLowerColor := func(a [2]uint8, b [2]uint8) bool {
		return (a[0]*13 + a[1]) < (b[0]*13 + b[1])
	}

	highestCount := 0
	teamColors := [2]uint8{0, 0}

	for colors, count := range colorCount {
		shouldSwap := (count > highestCount) || (count == highestCount && isLowerColor(colors, teamColors))

		if shouldSwap {
			teamColors = colors
			highestCount = count
		}
	}

	return teamColors
}

func (t Team) String() string {
	playerCount := len(t.Players)

	if 0 == playerCount || playerCount > 4 {
		return t.Name.ToPlainString()
	}

	playerNames := make([]string, 0)

	for _, p := range t.Players {
		playerNames = append(playerNames, p.Name.ToPlainString())
	}

	playerNames = qutil.StripQuakeFixes(playerNames)

	sort.Slice(playerNames, func(i, j int) bool {
		return strings.ToLower(playerNames[i]) < strings.ToLower(playerNames[j])
	})

	return fmt.Sprintf("%s (%s)", t.Name.ToPlainString(), strings.Join(playerNames, ", "))
}

func FromPlayers(players []qclient.Client) []Team {
	playersPerTeamName := make(map[string][]qclient.Client, 0)
	teamNamePerId := make(map[string]qstring.QuakeString, 0)

	for _, player := range players {
		teamId := player.Team.ToPlainString()
		playersPerTeamName[teamId] = append(playersPerTeamName[teamId], player)
		teamNamePerId[teamId] = player.Team
	}

	teams := make([]Team, 0)

	for teamId, teamName := range teamNamePerId {
		teams = append(teams, Team{
			Name:    teamName,
			Players: playersPerTeamName[teamId],
		})
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name.ToPlainString() < teams[j].Name.ToPlainString()
	})

	return teams
}
