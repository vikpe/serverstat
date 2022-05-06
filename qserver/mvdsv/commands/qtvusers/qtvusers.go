package qtvusers

import (
	"errors"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

func Send(udpClient udpclient.UdpClient, address string) ([]qclient.Client, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
	return ParseResponse(udpClient.Request(address, statusPacket, expectedHeader))
}

func ParseResponse(responseBody []byte, err error) ([]qclient.Client, error) {
	if err != nil {
		return make([]qclient.Client, 0), err
	}

	// example response body: 12 "djevulsk" "serp" "player" "rst" "twitch.tv/vikpe"
	fullText := string(responseBody)
	const QuoteChar = "\""

	if !strings.Contains(fullText, QuoteChar) {
		return make([]qclient.Client, 0), errors.New("invalid response body")
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

	return clients, nil
}
