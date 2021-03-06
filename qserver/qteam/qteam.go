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
	Name    qstring.QuakeString `json:"name"`
	Players []qclient.Client    `json:"players"`
}

type TeamExport struct {
	Name      qstring.QuakeString `json:"name"`
	NameColor string              `json:"name_color"`
	Frags     int                 `json:"frags"`
	Colors    [2]uint8            `json:"colors"`
	Players   []qclient.Client    `json:"players"`
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
	if 0 == len(players) {
		return make([]Team, 0)
	} else if 1 == len(players) {
		return []Team{{
			Name:    players[0].Team,
			Players: []qclient.Client{players[0]},
		}}
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Team.ToPlainString() < players[j].Team.ToPlainString()
	})

	teams := make([]Team, 0)
	currentTeamIndex := -1
	currentTeamName := "____________________"

	for _, player := range players {
		playerTeamName := player.Team.ToPlainString()

		if currentTeamName != playerTeamName {
			teams = append(teams, Team{
				Name:    player.Team,
				Players: []qclient.Client{player},
			})
			currentTeamIndex++
			currentTeamName = playerTeamName
		} else {
			teams[currentTeamIndex].Players = append(teams[currentTeamIndex].Players, player)
		}
	}

	for _, team := range teams {
		qclient.SortPlayers(team.Players)
	}

	return teams
}
