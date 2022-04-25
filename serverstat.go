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
	response, err := UdpRequest(address, statusPacket, expectedHeader)

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

	for scanner.Scan() {
		var client, err = parseClientString(scanner.Text())

		if err != nil {
			continue
		}

		if client.IsSpec {
			qserver.Spectators = append(qserver.Spectators, Spectator{
				Name:    client.Name,
				NameInt: client.NameInt,
				IsBot:   client.IsBot,
			})
		} else {
			qserver.Players = append(qserver.Players, client.Player)
		}
	}

	qserver.Address = address
	qserver.Title = qserver.Settings["hostname"]
	qserver.NumPlayers = len(qserver.Players)
	qserver.NumSpectators = len(qserver.Spectators)

	qtvServer, _ := StatQtv(address)
	qserver.QtvAddress = qtvServer.Address

	return qserver, nil
}

func StatQtv(address string) (QtvServer, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	response, err := UdpRequest(address, statusPacket, expectedHeader)

	if err != nil {
		return QtvServer{}, err
	}

	responseBody := response[5:]
	reader := csv.NewReader(strings.NewReader(string(responseBody)))
	reader.Comma = ' '

	record, err := reader.Read()
	if err != nil {
		return QtvServer{}, err
	}

	const (
		IndexTitle   = 2
		IndexAddress = 3
	)

	if record[IndexAddress] == "" {
		// these are the servers that are not configured correctly,
		// that means they are not reporting their qtv ip as they should.
		return QtvServer{}, err
	}

	return QtvServer{
		Title:      record[IndexTitle],
		Address:    record[IndexAddress],
		Spectators: make([]string, 0),
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
