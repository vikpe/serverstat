package serverstat

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseClientRecord(clientRecord []string) (Client, error) {
	columnCount := len(clientRecord)
	const ExpectedColumnCount = 9

	if columnCount != ExpectedColumnCount {
		err := errors.New(fmt.Sprintf("invalid player column count %d.", columnCount))
		return Client{}, err
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

	name := quakeTextToPlainText(nameQuakeStr)
	nameInt := stringToIntArray(name)
	team := quakeTextToPlainText(clientRecord[IndexTeam])
	teamInt := stringToIntArray(team)
	colorTop := stringToInt(clientRecord[IndexColorTop])
	colorBottom := stringToInt(clientRecord[IndexColorBottom])
	ping := stringToInt(clientRecord[IndexPing])

	return Client{
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
