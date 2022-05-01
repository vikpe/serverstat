package serverstat

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/vikpe/qw-serverstat/quakeserver"
	"github.com/vikpe/qw-serverstat/quakeserver/qtvstream"
	"github.com/vikpe/qw-serverstat/quaketext/qstring"
	"github.com/vikpe/udpclient"
)

func GetServerInfo(address string) (quakeserver.QuakeServer, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}

	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return quakeserver.QuakeServer{}, err
	}

	responseBody := response[len(expectedHeader):]
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))
	scanner.Scan()

	qserver := quakeserver.New()
	qserver.Settings = parseSettingsString(scanner.Text())

	if val, ok := qserver.Settings["hostname"]; ok {
		qserver.Settings["hostname"] = qstring.ToPlainString(val)
	}
	if val, ok := qserver.Settings["map"]; ok {
		qserver.Map = val
	}
	if val, ok := qserver.Settings["maxclients"]; ok {
		value, _ := strconv.Atoi(val)
		qserver.MaxPlayers = uint8(value)
	}
	if val, ok := qserver.Settings["maxspectators"]; ok {
		value, _ := strconv.Atoi(val)
		qserver.MaxSpectators = uint8(value)
	}

	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	players, spectators := parseClientsStrings(clientStrings)
	qserver.Players = players
	qserver.Spectators = spectators

	qserver.Address = address
	qserver.Title = qserver.Settings["hostname"]
	qserver.NumPlayers = uint8(len(qserver.Players))
	qserver.NumSpectators = uint8(len(qserver.Spectators))

	qtvServerStream, _ := GetQtvStreamInfo(address)
	qserver.QtvStream = qtvServerStream

	return qserver, nil
}

func parseSettingsString(settingsString string) map[string]string {
	settingsLines := strings.FieldsFunc(settingsString, func(r rune) bool { return r == '\\' })
	settings := make(map[string]string, len(settingsLines))

	for i := 0; i < len(settingsLines)-1; i += 2 {
		settings[settingsLines[i]] = settingsLines[i+1]
	}

	return settings
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

	name := qstring.ToPlainString(nameQuakeStr)
	team := qstring.ToPlainString(clientRecord[IndexTeam])
	colorTop := stringToInt(clientRecord[IndexColorTop])
	colorBottom := stringToInt(clientRecord[IndexColorBottom])
	ping := stringToInt(clientRecord[IndexPing])

	return quakeserver.Client{
		Player: quakeserver.Player{
			Name:    name,
			NameRaw: []byte(clientRecord[IndexName]),
			Team:    team,
			TeamRaw: []byte(clientRecord[IndexTeam]),
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
	return valueAsInt
}

func GetQtvUsers(address string) []string {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
	udpClient := udpclient.New()
	response, _ := udpClient.Request(address, statusPacket, expectedHeader)
	responseBody := response[len(expectedHeader):]
	return parseQtvusersResponseBody(responseBody)
}

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
	namesText = qstring.ToPlainString(namesText)

	return strings.Split(namesText, "\" \"")
}

func GetQtvStreamInfo(address string) (qtvstream.QtvStream, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return qtvstream.QtvStream{}, err
	}

	responseBody := response[5:]
	reader := csv.NewReader(strings.NewReader(string(responseBody)))
	reader.Comma = ' '

	record, err := reader.Read()
	if err != nil {
		return qtvstream.QtvStream{}, err
	}

	const (
		IndexTitle       = 2
		IndexAddress     = 3
		IndexClientCount = 4
	)

	if record[IndexAddress] == "" {
		// these are the servers that are not configured correctly,
		// that means they are not reporting their qtv ip as they should.
		return qtvstream.QtvStream{}, err
	}

	numberOfSpectators := stringToInt(record[IndexClientCount])

	var spectatorNames []string

	if numberOfSpectators > 0 {
		spectatorNames = GetQtvUsers(address)
	} else {
		spectatorNames = make([]string, 0)
	}

	return qtvstream.QtvStream{
		Title:          record[IndexTitle],
		Url:            record[IndexAddress],
		SpectatorNames: spectatorNames,
		NumSpectators:  uint8(numberOfSpectators),
	}, nil
}

func GetServerInfoFromMany(addresses []string) []quakeserver.QuakeServer {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	servers := make([]quakeserver.QuakeServer, 0)

	for _, address := range addresses {
		wg.Add(1)

		go func(address string) {
			defer wg.Done()

			qserver, err := GetServerInfo(address)

			if err != nil {
				return
			}

			mutex.Lock()
			servers = append(servers, qserver)
			mutex.Unlock()
		}(address)
	}

	wg.Wait()

	return servers
}
