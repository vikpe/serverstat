package qclient

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vikpe/serverstat/quaketext/qstring"
)

type Client struct {
	Name    string
	NameRaw []rune
	Team    string
	TeamRaw []rune
	Skin    string
	Colors  [2]uint8
	Frags   int
	Ping    int
	Time    uint8
	IsBot   bool
}

func FromStrings(clientStrings []string) []Client {
	clients := make([]Client, 0)

	for _, clientStr := range clientStrings {
		client, err := NewFromString(clientStr)

		if err != nil {
			continue
		}

		clients = append(clients, client)
	}

	return clients
}

func NewFromString(clientString string) (Client, error) {
	reader := csv.NewReader(strings.NewReader(clientString))
	reader.Comma = ' '

	clientRecord, err := reader.Read()
	if err != nil {
		return Client{}, err
	}

	client, err := NewFromRecord(clientRecord)

	if err != nil {
		return Client{}, err
	}

	return client, nil
}

func NewFromRecord(clientRecord []string) (Client, error) {
	minimumColumnCount := uint8(8)
	columnCount := uint8(len(clientRecord))

	if columnCount < minimumColumnCount {
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
	nameQuakeStr = strings.TrimPrefix(nameQuakeStr, SpectatorPrefix)

	name := qstring.ToPlainString(nameQuakeStr)
	frags := StringToInt(clientRecord[IndexFrags])
	colorTop := StringToInt(clientRecord[IndexColorTop])
	colorBottom := StringToInt(clientRecord[IndexColorBottom])
	ping := StringToInt(clientRecord[IndexPing])

	team := ""
	teamRaw := make([]rune, 0)

	if columnCount-1 >= IndexTeam {
		team = qstring.ToPlainString(clientRecord[IndexTeam])
		teamRaw = []rune(clientRecord[IndexTeam])
	}

	return Client{
		Name:    name,
		NameRaw: []rune(nameQuakeStr),
		Team:    team,
		TeamRaw: teamRaw,
		Skin:    clientRecord[IndexSkin],
		Colors:  [2]uint8{uint8(colorTop), uint8(colorBottom)},
		Frags:   frags,
		Ping:    ping,
		Time:    uint8(StringToInt(clientRecord[IndexTime])),
		IsBot:   isBotName(name) || isBotPing(ping),
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

func StringToInt(value string) int {
	valueAsInt, _ := strconv.Atoi(value)
	return valueAsInt
}
