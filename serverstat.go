package serverstat

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/vikpe/qw-serverstat/qserver"
	"github.com/vikpe/qw-serverstat/quaketext/qstring"
	"github.com/vikpe/udpclient"
)

func GetServerInfo(address string) (qserver.GenericServer, error) {
	// query
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}

	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return qserver.GenericServer{}, err
	}

	// resulting server
	server := qserver.NewGenericServer()

	// settings
	responseBody := response[len(expectedHeader):]
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))
	scanner.Scan()
	server.Settings = parseSettingsString(scanner.Text())

	/*if val, ok := server.Settings["hostname"]; ok {
		server.Settings["hostname"] = qstring.ToPlainString(val)
	}
	if val, ok := server.Settings["map"]; ok {
		server.Map = val
	}
	if val, ok := server.Settings["maxclients"]; ok {
		value, _ := strconv.Atoi(val)
		server.MaxPlayers = uint8(value)
	}
	if val, ok := server.Settings["maxspectators"]; ok {
		value, _ := strconv.Atoi(val)
		server.MaxSpectators = uint8(value)
	}*/

	//server.Title = server.Settings["hostname"]
	//server.NumPlayers = uint8(len(server.Players))
	//server.NumSpectators = uint8(len(server.Spectators))

	// clients
	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	if len(clientStrings) > 0 {
		clientColumnCount := uint8(0)

		if qserver.IsGameServer(server) {
			clientColumnCount = 9

		} else if qserver.IsQtvServer(server) || qserver.IsProxyServer(server) {
			clientColumnCount = 8
		}

		if clientColumnCount > 0 {
			server.Clients = parseClientStrings(clientStrings, clientColumnCount)
		}
	}

	// extra
	if qserver.IsGameServer(server) {
		qtvServerStream, _ := GetQtvStreamInfo(address)
		server.ExtraInfo.QtvStream = qtvServerStream
	}

	// common
	server.Address = address
	server.NumClients = uint8(len(server.Clients))

	return server, nil
}

func parseSettingsString(settingsString string) map[string]string {
	settingsLines := strings.FieldsFunc(settingsString, func(r rune) bool { return r == '\\' })
	settings := make(map[string]string, len(settingsLines))

	for i := 0; i < len(settingsLines)-1; i += 2 {
		settings[settingsLines[i]] = settingsLines[i+1]
	}

	return settings
}

func parseClientStrings(clientStrings []string, expectedColumnCount uint8) []qserver.Client {
	clients := make([]qserver.Client, 0)

	for _, clientStr := range clientStrings {
		clientRecord, err := clientStringToRecord(clientStr)

		if err != nil {
			continue
		}

		client, err := parseClientRecord(clientRecord, expectedColumnCount)

		if err != nil {
			continue
		}

		clients = append(clients, client)
	}

	return clients
}

func clientStringToRecord(clientStr string) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(clientStr))
	reader.Comma = ' '

	clientRecord, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return clientRecord, nil
}

func parseClientRecord(clientRecord []string, expectedColumnCount uint8) (qserver.Client, error) {
	columnCount := uint8(len(clientRecord))

	if columnCount != expectedColumnCount {
		err := errors.New(fmt.Sprintf("invalid player column count %d.", columnCount))
		return qserver.Client{}, err
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
	frags := stringToInt(clientRecord[IndexFrags])
	colorTop := stringToInt(clientRecord[IndexColorTop])
	colorBottom := stringToInt(clientRecord[IndexColorBottom])
	ping := stringToInt(clientRecord[IndexPing])

	team := ""
	teamRaw := make([]rune, 0)

	if columnCount-1 >= IndexTeam {
		team = qstring.ToPlainString(clientRecord[IndexTeam])
		teamRaw = []rune(clientRecord[IndexTeam])
	}

	return qserver.Client{
		Name:    name,
		NameRaw: []rune(nameQuakeStr),
		Team:    team,
		TeamRaw: teamRaw,
		Skin:    clientRecord[IndexSkin],
		Colors:  [2]uint8{uint8(colorTop), uint8(colorBottom)},
		Frags:   frags,
		Ping:    ping,
		Time:    uint8(stringToInt(clientRecord[IndexTime])),
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

func stringToInt(value string) int {
	valueAsInt, _ := strconv.Atoi(value)
	return valueAsInt
}

func GetQtvUsers(address string) ([]string, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return nil, err
	}

	responseBody := response[len(expectedHeader):]
	return parseQtvusersResponseBody(responseBody), nil
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

func GetQtvStreamInfo(address string) (qserver.QtvStream, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return qserver.QtvStream{}, err
	}

	responseBody := response[5:]
	reader := csv.NewReader(strings.NewReader(string(responseBody)))
	reader.Comma = ' '

	record, err := reader.Read()
	if err != nil {
		return qserver.QtvStream{}, err
	}

	const (
		IndexTitle       = 2
		IndexAddress     = 3
		IndexClientCount = 4
	)

	if record[IndexAddress] == "" {
		// these are the servers that are not configured correctly,
		// that means they are not reporting their qtv ip as they should.
		return qserver.QtvStream{}, err
	}

	numberOfSpectators := stringToInt(record[IndexClientCount])

	var spectatorNames []string

	if numberOfSpectators > 0 {
		spectatorNames, err = GetQtvUsers(address)

		if err != nil {
			spectatorNames = make([]string, 0)
		}
	} else {
		spectatorNames = make([]string, 0)
	}

	return qserver.QtvStream{
		Title:          record[IndexTitle],
		Url:            record[IndexAddress],
		SpectatorNames: spectatorNames,
		NumSpectators:  uint8(numberOfSpectators),
	}, nil
}

func GetServerInfoFromMany(addresses []string) []qserver.GenericServer {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	servers := make([]qserver.GenericServer, 0)

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
