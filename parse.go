package serverstat

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vikpe/qw-serverstat/quaketext"
)

func parseQtvusersResponseBody(responseBody []byte) []string {
	// example response body: 12 "djevulsk" "serp" "player" "rst" "twitch.tv/vikpe"
	fullText := string(responseBody)
	const QuoteChar = "\""

	if !strings.Contains(fullText, QuoteChar) {
		return make([]string, 0)
	}

	indexFirstQuote := strings.Index(fullText, QuoteChar)
	indexLastQuote := strings.LastIndex(fullText, QuoteChar)
	namesText := fullText[indexFirstQuote+1 : indexLastQuote]
	namesText = quaketext.ToPlainText(namesText)

	return strings.Split(namesText, "\" \"")
}

func parseClientRecord(clientRecord []string) (client, error) {
	columnCount := len(clientRecord)
	const ExpectedColumnCount = 9

	if columnCount != ExpectedColumnCount {
		err := errors.New(fmt.Sprintf("invalid player column count %d.", columnCount))
		return client{}, err
	}

	const (
		IndexFrags              = 1
		IndexTime               = 2
		IndexPing               = 3
		IndexName               = 4
		IndexSkin               = 5
		IndexColorTop           = 6
		IndexColorBottom        = 7
		IndexTeam               = 8
		SpectatorPrefix  string = "\\s\\"
	)

	nameQuakeStr := clientRecord[IndexName]

	isSpec := strings.HasPrefix(nameQuakeStr, SpectatorPrefix)
	if isSpec {
		nameQuakeStr = strings.TrimPrefix(nameQuakeStr, SpectatorPrefix)
	}

	name := quaketext.ToPlainText(nameQuakeStr)
	nameInt := stringToIntArray(name)
	team := quaketext.ToPlainText(clientRecord[IndexTeam])
	teamInt := stringToIntArray(team)
	colorTop := stringToInt(clientRecord[IndexColorTop])
	colorBottom := stringToInt(clientRecord[IndexColorBottom])
	ping := stringToInt(clientRecord[IndexPing])

	return client{
		Player: Player{
			Name:    name,
			NameInt: nameInt,
			Team:    team,
			TeamInt: teamInt,
			Skin:    clientRecord[IndexSkin],
			Colors:  [2]int{colorTop, colorBottom},
			Frags:   stringToInt(clientRecord[IndexFrags]),
			Ping:    ping,
			Time:    stringToInt(clientRecord[IndexTime]),
			IsBot:   isBotName(name) || isBotPing(ping),
		},
		IsSpec: isSpec,
	}, nil

}

func parseClientString(clientStr string) (client, error) {
	reader := csv.NewReader(strings.NewReader(clientStr))
	reader.Comma = ' '

	clientRecord, err := reader.Read()
	if err != nil {
		return client{}, nil
	}

	return parseClientRecord(clientRecord)
}

func parseClientsStrings(clientStrings []string) ([]Player, []Spectator) {
	players := make([]Player, 0)
	spectators := make([]Spectator, 0)

	for _, clientStr := range clientStrings {
		var client, err = parseClientString(clientStr)

		if err != nil {
			continue
		}

		if client.IsSpec {
			spectators = append(spectators, Spectator{
				Name:    client.Name,
				NameInt: client.NameInt,
				IsBot:   client.IsBot,
			})
		} else {
			players = append(players, client.Player)
		}
	}

	return players, spectators
}

func isBotName(name string) bool {
	switch name {
	case
		"[ServeMe]",
		"twitch.tv/vikpe":
		return true
	}
	return false
}

func isBotPing(ping int) bool {
	switch ping {
	case
		10:
		return true
	}
	return false
}

func stringToInt(value string) int {
	valueAsInt, _ := strconv.Atoi(value)
	return valueAsInt
}

func stringToIntArray(value string) []int {
	intArr := make([]int, len(value))

	for i := range value {
		intArr[i] = int(value[i])
	}

	return intArr
}
