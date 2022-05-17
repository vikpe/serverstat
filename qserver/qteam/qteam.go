package qteam

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
)

type Team struct {
	Name    qstring.QuakeString
	Players []qclient.Client
}

type TeamExport struct {
	Name    qstring.QuakeString
	Frags   int
	Colors  [2]uint8
	Players []qclient.Client
}

func Export(t Team) TeamExport {
	return TeamExport{
		Name:    t.Name,
		Colors:  t.Colors(),
		Frags:   t.Frags(),
		Players: t.Players,
	}
}

func (t Team) MarshalJSON() ([]byte, error) {
	return json.Marshal(Export(t))
}

func (t Team) Frags() int {
	frags := 0

	for _, p := range t.Players {
		frags += p.Frags
	}

	return frags
}

func (t Team) Colors() [2]uint8 {
	colorCount := make(map[[2]uint8]int, 0)

	for _, p := range t.Players {
		colorCount[p.Colors]++
	}

	highestCount := 0
	teamColors := [2]uint8{0, 0}

	for colorCombination, count := range colorCount {
		if count > highestCount {
			teamColors = colorCombination
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

	return teams
}
