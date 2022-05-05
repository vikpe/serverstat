package mvdsv

import (
	"encoding/csv"
	"encoding/json"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/udpclient"
)

type QtvStream struct {
	Title   string
	Url     string
	Clients []qclient.Client
}

func (q QtvStream) MarshalJSON() ([]byte, error) {
	if "" == q.Url {
		return json.Marshal("")
	} else {
		type QtvStreamJson struct {
			Title      string
			Url        string
			Clients    []qclient.Client
			NumClients uint8
		}

		return json.Marshal(QtvStreamJson{
			Title:      q.Title,
			Url:        q.Url,
			Clients:    q.Clients,
			NumClients: uint8(len(q.Clients)),
		})
	}
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

func GetQtvStreamInfo(address string) (QtvStream, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

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

	numberOfClients := qutil.StringToInt(record[IndexClientCount])

	var clients []qclient.Client

	if numberOfClients > 0 {
		clients, err = GetQtvusers(address)

		if err != nil {
			clients = make([]qclient.Client, 0)
		}
	} else {
		clients = make([]qclient.Client, 0)
	}

	return QtvStream{
		Title:   record[IndexTitle],
		Url:     record[IndexAddress],
		Clients: clients,
	}, nil
}
