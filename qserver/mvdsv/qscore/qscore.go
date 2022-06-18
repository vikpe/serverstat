package qscore

import (
	"fmt"
	"math"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
)

func FromModeAndPlayers(mode string, players []qclient.Client) int {
	botCount := getBotCount(players)
	playerCount := len(players)

	if botCount == playerCount {
		return 0
	}

	botPercentage := float64(botCount) / float64(playerCount)
	score := FromModeAndPlayerNames(mode, qclient.ClientNames(players))
	weightedScore := math.Round((1.0 - botPercentage) * float64(score))

	return int(weightedScore)
}

func FromModeAndPlayerNames(mode string, playerNames []string) int {
	playerCount := len(playerNames)

	if 0 == playerCount {
		return 0
	}

	maxScore := getMaxScoreByMode(mode)
	expectedPlayerCount := getExpectedPlayerCountByMode(mode)

	if playerCount < expectedPlayerCount {
		missingPlayers := expectedPlayerCount - playerCount

		if missingPlayers >= 2 {
			maxScore = 5 * float64(playerCount)
			//fmt.Println(mode, playerCount, maxScore)
		} else {
			fillPercent := float64(playerCount) / float64(expectedPlayerCount)
			reductionFactor := fillPercent * fillPercent * fillPercent
			maxScore *= reductionFactor
			//fmt.Println(mode, reductionFactor)
		}
	}

	playerFactor := 1 / getAverageDiv(playerNames)
	score := int(playerFactor * maxScore)

	fmt.Println(score, " .. ", mode, strings.Join(playerNames, ", "), " factor ", playerFactor)

	return score
}

func getBotCount(clients []qclient.Client) int {
	count := 0

	for _, c := range clients {
		if c.IsBot() {
			count++
		}
	}

	return count
}

func getMaxScoreByMode(mode string) float64 {
	switch mode {
	case "1on1":
		return 50
	case "2on2":
		return 70
	case "4on4":
		return 130
	case "ffa":
		return 35
	case "coop":
		return 20
	case "race":
		return 15
	case "3on3":
		return 75
	}

	return 5
}

func getExpectedPlayerCountByMode(mode string) int {
	switch mode {
	case "1on1":
		return 2
	case "2on2":
		return 4
	case "3on3":
		return 6
	case "4on4":
		return 8
	case "ffa":
		return 2
	}

	return 0
}

func getPlayerDiv(name string) float64 {
	strippedName := stripName(name)

	if div, ok := PlayerDivs[strippedName]; ok {
		return div
	}
	const unknownPlayerDiv = 2.5
	return unknownPlayerDiv
}

func getAverageDiv(playerNames []string) float64 {
	totalScore := float64(0)

	for _, name := range playerNames {
		totalScore += getPlayerDiv(name)
	}

	return totalScore / float64(len(playerNames))
}

func stripName(name string) string {
	var result string

	if strings.ContainsRune(name, '•') {
		result = strings.ReplaceAll(name, "•", " ")
	} else {
		result = name
	}

	return strings.ToLower(strings.TrimSpace(result))
}
