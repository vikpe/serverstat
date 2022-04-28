package serverstat

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/vikpe/qw-serverstat/quaketext"
	"github.com/vikpe/qw-serverstat/qwnet"
)

type Player struct {
	Name    string
	NameRaw []uint16
	Team    string
	TeamRaw []uint16
	Skin    string
	Colors  [2]int
	Frags   int
	Ping    int
	Time    int
	IsBot   bool
}

type client struct {
	Player
	IsSpec bool
}

type Spectator struct {
	Name    string
	NameRaw []uint16
	IsBot   bool
}

type QtvStream struct {
	Title          string
	Url            string
	SpectatorNames []string
	NumSpectators  int
}

func (node *QtvStream) MarshalJSON() ([]byte, error) {
	if "" == node.Url {
		return json.Marshal(nil)
	} else {
		return json.Marshal(QtvStream{
			Title:          node.Title,
			Url:            node.Url,
			SpectatorNames: node.SpectatorNames,
			NumSpectators:  node.NumSpectators,
		})
	}
}

func newQtvStream() QtvStream {
	return QtvStream{
		Title:         "",
		Url:           "",
		NumSpectators: 0,
	}
}

type QuakeServer struct {
	Title         string
	Address       string
	QtvStream     QtvStream
	Map           string
	NumPlayers    int
	MaxPlayers    int
	NumSpectators int
	MaxSpectators int
	Players       []Player
	Spectators    []Spectator
	Settings      map[string]string
}

func newQuakeServer() QuakeServer {
	return QuakeServer{
		Title:      "",
		Address:    "",
		Settings:   map[string]string{},
		Players:    make([]Player, 0),
		Spectators: make([]Spectator, 0),
		QtvStream:  newQtvStream(),
	}
}

func Stat(address string) (QuakeServer, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}
	response, err := qwnet.UdpRequest(address, statusPacket, expectedHeader)

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
		qserver.Settings["hostname"] = quaketext.StringToPlainString(val)
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

	qtvServerStream, _ := statQtvStream(address)
	qserver.QtvStream = qtvServerStream

	return qserver, nil
}

func statQtvStreamUsers(address string) []string {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
	response, _ := qwnet.UdpRequest(address, statusPacket, expectedHeader)
	responseBody := response[len(expectedHeader):]
	return parseQtvusersResponseBody(responseBody)
}

func statQtvStream(address string) (QtvStream, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	response, err := qwnet.UdpRequest(address, statusPacket, expectedHeader)

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
		IndexTitle       = 2
		IndexAddress     = 3
		IndexClientCount = 4
	)

	if record[IndexAddress] == "" {
		// these are the servers that are not configured correctly,
		// that means they are not reporting their qtv ip as they should.
		return QtvStream{}, err
	}

	numberOfSpectators := stringToInt(record[IndexClientCount])

	var spectatorNames []string

	if numberOfSpectators > 0 {
		spectatorNames = statQtvStreamUsers(address)
	} else {
		spectatorNames = make([]string, 0)
	}

	return QtvStream{
		Title:          record[IndexTitle],
		Url:            record[IndexAddress],
		SpectatorNames: spectatorNames,
		NumSpectators:  numberOfSpectators,
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
