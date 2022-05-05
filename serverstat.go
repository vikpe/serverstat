package serverstat

import (
	"bufio"
	"encoding/csv"
	"strings"
	"sync"

	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtvstream"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

func GetServerInfo(address string) (qserver.GenericServer, error) {
	// request
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}

	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return qserver.GenericServer{}, err
	}

	// response
	responseBody := response[len(expectedHeader):]
	settingsString, clientStrings := parseServerInfoResponseBody(responseBody)

	// resulting server
	server := qserver.NewGenericServer()
	server.Address = address
	server.Settings = qsettings.New(settingsString)
	server.Clients = qclient.NewFromStrings(clientStrings)
	server.NumClients = uint8(len(server.Clients))

	// extra info
	version := qversion.New(server.Settings["*version"])

	if version.IsMvdsv() {
		qtvServerStream, _ := GetQtvStreamInfo(address)
		server.ExtraInfo.QtvStream = qtvServerStream
	}

	return server, nil
}

func parseServerInfoResponseBody(responseBody []byte) (string, []string) {
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))

	scanner.Scan()
	settingsString := scanner.Text()

	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	return settingsString, clientStrings
}

func GetQtvusers(address string) ([]qclient.Client, error) {
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

func parseQtvusersResponseBody(responseBody []byte) []qclient.Client {
	// example response body: 12 "djevulsk" "serp" "player" "rst" "twitch.tv/vikpe"
	fullText := string(responseBody)
	const QuoteChar = "\""

	if !strings.Contains(fullText, QuoteChar) {
		return make([]qclient.Client, 0)
	}

	indexFirstQuote := strings.Index(fullText, QuoteChar)
	indexLastQuote := strings.LastIndex(fullText, QuoteChar)
	namesText := fullText[indexFirstQuote+1 : indexLastQuote]

	clients := make([]qclient.Client, 0)
	namesRaw := strings.Split(namesText, "\" \"")

	for _, nameRaw := range namesRaw {
		clients = append(clients, qclient.Client{
			Name:    qstring.ToPlainString(nameRaw),
			NameRaw: []rune(nameRaw),
		})
	}

	return clients
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

	numberOfClients := qclient.StringToInt(record[IndexClientCount])

	var clients []qclient.Client

	if numberOfClients > 0 {
		clients, err = GetQtvusers(address)

		if err != nil {
			clients = make([]qclient.Client, 0)
		}
	} else {
		clients = make([]qclient.Client, 0)
	}

	return qtvstream.QtvStream{
		Title:      record[IndexTitle],
		Url:        record[IndexAddress],
		Clients:    clients,
		NumClients: uint8(numberOfClients),
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

			server, err := GetServerInfo(address)

			if err != nil {
				return
			}

			mutex.Lock()
			servers = append(servers, server)
			mutex.Unlock()
		}(address)
	}

	wg.Wait()

	return servers
}
