package qtvusers

import (
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

func New(address string) ([]qclient.Client, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return nil, err
	}

	responseBody := response[len(expectedHeader):]
	return ParseResponseBody(responseBody), nil
}

func ParseResponseBody(responseBody []byte) []qclient.Client {
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
