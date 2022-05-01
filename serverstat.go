package serverstat

import (
	"bufio"
	"encoding/csv"
	"strconv"
	"strings"
	"sync"

	"github.com/vikpe/qw-serverstat/qtvstream"
	"github.com/vikpe/qw-serverstat/quakeserver"
	"github.com/vikpe/qw-serverstat/quaketext"
	"github.com/vikpe/udpclient"
)

func Stat(address string) (quakeserver.QuakeServer, error) {
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

	settings := strings.FieldsFunc(scanner.Text(), func(r rune) bool { return r == '\\' })

	qserver := quakeserver.New()

	for i := 0; i < len(settings)-1; i += 2 {
		qserver.Settings[settings[i]] = settings[i+1]
	}

	if val, ok := qserver.Settings["hostname"]; ok {
		qserver.Settings["hostname"] = quaketext.NewFromString(val).ToPlainString()
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

	qtvServerStream, _ := statQtvStream(address)
	qserver.QtvStream = qtvServerStream

	return qserver, nil
}

func statQtvStreamUsers(address string) []string {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
	udpClient := udpclient.New()
	response, _ := udpClient.Request(address, statusPacket, expectedHeader)
	responseBody := response[len(expectedHeader):]
	return parseQtvusersResponseBody(responseBody)
}

func statQtvStream(address string) (qtvstream.QtvStream, error) {
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
		spectatorNames = statQtvStreamUsers(address)
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

func StatMany(addresses []string) []quakeserver.QuakeServer {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	servers := make([]quakeserver.QuakeServer, 0)

	for _, address := range addresses {
		wg.Add(1)

		go func(address string) {
			defer wg.Done()

			qserver, err := Stat(address)

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
