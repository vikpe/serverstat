package serverstat

import (
	"bufio"
	"encoding/csv"
	"strconv"
	"strings"
	"sync"
)

func Stat(address string) (QuakeServer, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}
	response, err := udpRequest(address, statusPacket, expectedHeader)

	if err != nil {
		return QuakeServer{}, err
	}

	responseBody := response[len(expectedHeader):]
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))
	scanner.Scan()

	settings := strings.FieldsFunc(scanner.Text(), func(r rune) bool { return r == '\\' })

	qserver := newQuakeServer()

	for i := 0; i < len(settings)-1; i += 2 {
		qserver.Settings[settings[i]] = settings[i+1]
	}

	if val, ok := qserver.Settings["hostname"]; ok {
		qserver.Settings["hostname"] = quakeTextToPlainText(val)

	}
	if val, ok := qserver.Settings["map"]; ok {
		qserver.Map = val
	}
	if val, ok := qserver.Settings["maxclients"]; ok {
		value, _ := strconv.Atoi(val)
		qserver.MaxPlayers = value
	}
	if val, ok := qserver.Settings["maxspectators"]; ok {
		value, _ := strconv.Atoi(val)
		qserver.MaxSpectators = value
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
	qserver.NumPlayers = len(qserver.Players)
	qserver.NumSpectators = len(qserver.Spectators)

	qtvServerStream, _ := StatServerQtvStream(address)
	qserver.QtvStream = qtvServerStream

	return qserver, nil
}

/*
// TODO: fixme
func StatServerQtvUsers(address string) []string {

	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x20}
	response, err := serverstat.UdpRequest("95.216.18.118:28001", statusPacket, expectedHeader)
	scanner := bufio.NewScanner(strings.NewReader(string(response)))
	scanner.Scan()
}
*/

func StatServerQtvStream(address string) (QtvStream, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	response, err := udpRequest(address, statusPacket, expectedHeader)

	if err != nil {
		return QtvStream{}, err
	}

	responseBody := response[5:]
	reader := csv.NewReader(strings.NewReader(string(responseBody)))
	reader.Comma = ' '

	record, err := reader.Read()
	if err != nil {
		return QtvStream{}, err
	}

	const (
		IndexId          = 1
		IndexTitle       = 2
		IndexAddress     = 3
		IndexClientCount = 4
	)

	if record[IndexAddress] == "" {
		// these are the servers that are not configured correctly,
		// that means they are not reporting their qtv ip as they should.
		return QtvStream{}, err
	}

	return QtvStream{
		Id:            stringToInt(record[IndexId]),
		Title:         record[IndexTitle],
		Url:           record[IndexAddress],
		NumSpectators: stringToInt(record[IndexClientCount]),
	}, nil
}

func StatMany(addresses []string) []QuakeServer {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	servers := make([]QuakeServer, 0)

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
