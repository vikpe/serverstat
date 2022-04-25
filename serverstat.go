package serverstat

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Stat(address string) (QuakeServer, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return QuakeServer{}, err
	}
	defer conn.Close()

	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	buffer := make([]byte, 8192)
	bufferLength := 0

	const (
		Retries     = 3
		TimeoutInMs = 500
	)

	for i := 0; i < Retries; i++ {
		conn.SetDeadline(timeInFuture(TimeoutInMs))

		_, err = conn.Write(statusPacket)
		if err != nil {
			return QuakeServer{}, err
		}

		conn.SetDeadline(timeInFuture(TimeoutInMs))
		bufferLength, err = conn.Read(buffer)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return QuakeServer{}, err
	}

	validHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}
	responseHeader := buffer[:len(validHeader)]
	isValidHeader := bytes.Equal(responseHeader, validHeader)
	if !isValidHeader {
		log.Println(address + ": Response error")
		return QuakeServer{}, err
	}

	responseBody := buffer[len(validHeader):bufferLength]
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))
	scanner.Scan()

	settings := strings.FieldsFunc(scanner.Text(), func(r rune) bool {
		if r == '\\' {
			return true
		}
		return false
	})

	qserver := newQuakeServer()
	qserver.Address = address

	for i := 0; i < len(settings)-1; i += 2 {
		qserver.Settings[settings[i]] = settings[i+1]
	}

	if val, ok := qserver.Settings["hostname"]; ok {
		qserver.Settings["hostname"] = quakeTextToPlainText(val)
		qserver.Title = qserver.Settings["hostname"]
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
		reader := csv.NewReader(strings.NewReader(scanner.Text()))
		reader.Comma = ' '

		clientRecord, err := reader.Read()
		if err != nil {
			continue
		}

		client, err := parseClientRecord(clientRecord)
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

	qserver.NumPlayers = len(qserver.Players)
	qserver.NumSpectators = len(qserver.Spectators)

	qtvServer, _ := StatQtv(address)
	qserver.QtvAddress = qtvServer.Address

	return qserver, nil
}

func StatQtv(address string) (QtvServer, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return QtvServer{}, err
	}
	defer conn.Close()

	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	buffer := make([]byte, 8192)
	bufferLength := 0

	const (
		Retries     = 3
		TimeoutInMs = 500
	)

	for i := 0; i < Retries; i++ {
		conn.SetDeadline(timeInFuture(TimeoutInMs))

		_, err = conn.Write(statusPacket)
		if err != nil {
			return QtvServer{}, err
		}

		conn.SetDeadline(timeInFuture(TimeoutInMs))
		bufferLength, err = conn.Read(buffer)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		// no logging here. it seems that servers may not reply if they do not support
		// this specific "32" status request.
		return QtvServer{}, err
	}

	validHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	responseHeader := buffer[:len(validHeader)]
	isValidHeader := bytes.Equal(responseHeader, validHeader)
	if !isValidHeader {
		// some servers react to the specific "32" status message but will send the regular
		// status message because they misunderstood our command.
		return QtvServer{}, err
	}

	responseBody := buffer[5:bufferLength]
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

func timeInFuture(delta int) time.Time {
	return time.Now().Add(time.Duration(delta) * time.Millisecond)
}
