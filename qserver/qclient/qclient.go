package qclient

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/serverstat/qutil"
)

type Client struct {
	Name   qstring.QuakeString
	Team   qstring.QuakeString
	Skin   string
	Colors [2]uint8
	Frags  int
	Ping   int
	Time   uint8
	CC     string
	IsBot  bool
}

func NewFromStrings(clientStrings []string) []Client {
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

	minimumColumnCount := uint8(8)
	columnCount := uint8(len(clientRecord))

	if columnCount < minimumColumnCount {
		err := errors.New(fmt.Sprintf("invalid client column count %d, expects at least %d", columnCount, minimumColumnCount))
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
		IndexFlag               = 9
		SpectatorPrefix  string = "\\s\\"
	)

	nameQuakeStr := clientRecord[IndexName]
	nameQuakeStr = strings.TrimPrefix(nameQuakeStr, SpectatorPrefix)

	name := qstring.New(nameQuakeStr)
	frags := qutil.StringToInt(clientRecord[IndexFrags])
	colorTop := qutil.StringToInt(clientRecord[IndexColorTop])
	colorBottom := qutil.StringToInt(clientRecord[IndexColorBottom])
	ping := qutil.StringToInt(clientRecord[IndexPing])

	var lastIndex = columnCount - 1

	team := qstring.New("")

	if lastIndex >= IndexTeam {
		team = qstring.New(clientRecord[IndexTeam])
	}

	flag := ""

	if lastIndex >= IndexFlag {
		flag = clientRecord[IndexFlag]
	}

	return Client{
		Name:   name,
		Team:   team,
		Skin:   clientRecord[IndexSkin],
		Colors: [2]uint8{uint8(colorTop), uint8(colorBottom)},
		Frags:  frags,
		Ping:   ping,
		Time:   uint8(qutil.StringToInt(clientRecord[IndexTime])),
		CC:     flag,
		IsBot:  IsBotName(name.ToPlainString()) || IsBotPing(ping),
	}, nil

}

func IsBotName(name string) bool {
	if 0 == len(name) {
		return false
	}

	knownBotNames := []string{
		"[ServeMe]",
		"twitch.tv/vikpe",
	}

	return strings.Contains(strings.Join(knownBotNames, "\""), name)
}

func IsBotPing(ping int) bool {
	return 10 == ping
}
