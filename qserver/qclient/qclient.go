package qclient

import (
	"encoding/csv"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/ssoroka/slice"
	"github.com/vikpe/serverstat/qserver/qclient/bot"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/serverstat/qutil"
)

type Client struct {
	Name   qstring.QuakeString `json:"name"`
	Team   qstring.QuakeString `json:"team"`
	Skin   string              `json:"skin"`
	Colors [2]uint8            `json:"colors"`
	Frags  int                 `json:"frags"`
	Ping   int                 `json:"ping"`
	Time   int                 `json:"time"`
	CC     string              `json:"cc"`
}

func (client Client) IsBot() bool {
	return bot.IsBotPing(client.Ping) || bot.IsBotName(client.Name.ToPlainString())
}

func (client Client) IsHuman() bool {
	return !client.IsBot()
}

func (client Client) IsSpectator() bool {
	return client.Ping < 0
}

func (client Client) IsPlayer() bool {
	return !client.IsSpectator()
}

func (client Client) MarshalJSON() ([]byte, error) {
	return qutil.MarshalNoEscapeHtml(Export(client))
}

type ClientExport struct {
	Name      qstring.QuakeString `json:"name"`
	NameColor string              `json:"name_color"`
	Team      qstring.QuakeString `json:"team"`
	TeamColor string              `json:"team_color"`
	Skin      string              `json:"skin"`
	Colors    [2]uint8            `json:"colors"`
	Frags     int                 `json:"frags"`
	Ping      int                 `json:"ping"`
	Time      int                 `json:"time"`
	CC        string              `json:"cc"`
	IsBot     bool                `json:"is_bot"`
}

func Export(client Client) ClientExport {
	return ClientExport{
		Name:      client.Name,
		NameColor: client.Name.ToColorCodes(),
		Team:      client.Team,
		TeamColor: client.Team.ToColorCodes(),
		Skin:      client.Skin,
		Colors:    client.Colors,
		Frags:     client.Frags,
		Ping:      client.Ping,
		Time:      client.Time,
		CC:        client.CC,
		IsBot:     client.IsBot(),
	}
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

	colorTop := qutil.StringToInt(clientRecord[IndexColorTop])
	colorBottom := qutil.StringToInt(clientRecord[IndexColorBottom])
	ping := qutil.StringToInt(clientRecord[IndexPing])

	var indexCount = columnCount - 1

	team := qstring.New("")

	if indexCount >= IndexTeam {
		team = qstring.New(clientRecord[IndexTeam])
	}

	flag := ""

	if indexCount >= IndexFlag {
		flag = clientRecord[IndexFlag]
	}

	return Client{
		Name:   qstring.New(nameQuakeStr),
		Team:   team,
		Skin:   clientRecord[IndexSkin],
		Colors: [2]uint8{uint8(colorTop), uint8(colorBottom)},
		Frags:  qutil.StringToInt(clientRecord[IndexFrags]),
		Ping:   ping,
		Time:   qutil.StringToInt(clientRecord[IndexTime]),
		CC:     flag,
	}, nil
}

func ClientNames(clients []Client) []string {
	if 0 == len(clients) {
		return make([]string, 0)
	}

	return slice.Map[Client, string](clients, func(client Client) string {
		return client.Name.ToPlainString()
	})
}

func SortPlayers(players []Client) {
	if len(players) < 2 {
		return
	}

	sort.Slice(players, func(i, j int) bool {
		if players[i].Frags == players[j].Frags {
			return strings.ToLower(players[i].Name.ToPlainString()) < strings.ToLower(players[j].Name.ToPlainString())
		} else {
			return players[i].Frags > players[j].Frags
		}
	})
}
