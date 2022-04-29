package serverstat

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vikpe/qw-serverstat/quakeserver"
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
	namesText = quaketext.StringToPlainString(namesText)

	return strings.Split(namesText, "\" \"")
}

func parseClientRecord(clientRecord []string) (quakeserver.Client, error) {
	columnCount := len(clientRecord)
	const ExpectedColumnCount = 9

	if columnCount != ExpectedColumnCount {
		err := errors.New(fmt.Sprintf("invalid player column count %d.", columnCount))
		return quakeserver.Client{}, err
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

	name := quaketext.StringToPlainString(nameQuakeStr)
	team := quaketext.StringToPlainString(clientRecord[IndexTeam])
	colorTop := stringToInt(clientRecord[IndexColorTop])
	colorBottom := stringToInt(clientRecord[IndexColorBottom])
	ping := stringToInt(clientRecord[IndexPing])

	return quakeserver.Client{
		Player: quakeserver.Player{
			Name:    name,
			NameRaw: nameToRaw(clientRecord[IndexName]),
			Team:    team,
			TeamRaw: nameToRaw(clientRecord[IndexTeam]),
			Skin:    clientRecord[IndexSkin],
			Colors:  [2]uint8{uint8(colorTop), uint8(colorBottom)},
			Frags:   uint16(stringToInt(clientRecord[IndexFrags])),
			Ping:    uint16(ping),
			Time:    uint8(stringToInt(clientRecord[IndexTime])),
			IsBot:   isBotName(name) || isBotPing(ping),
		},
		IsSpec: isSpec,
	}, nil

}

func parseClientString(clientStr string) (quakeserver.Client, error) {
	reader := csv.NewReader(strings.NewReader(clientStr))
	reader.Comma = ' '

	clientRecord, err := reader.Read()
	if err != nil {
		return quakeserver.Client{}, nil
	}

	return parseClientRecord(clientRecord)
}

func parseClientsStrings(clientStrings []string) ([]quakeserver.Player, []quakeserver.Spectator) {
	players := make([]quakeserver.Player, 0)
	spectators := make([]quakeserver.Spectator, 0)

	for _, clientStr := range clientStrings {
		var client, err = parseClientString(clientStr)

		if err != nil {
			continue
		}

		if client.IsSpec {
			spectators = append(spectators, quakeserver.Spectator{
				Name:    client.Name,
				NameRaw: client.NameRaw,
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
	return int(valueAsInt)
}

func nameToRaw(value string) []uint16 {
	intArr := make([]uint16, len(value))

	for i := range value {
		intArr[i] = uint16(value[i])
	}

	return intArr
}
