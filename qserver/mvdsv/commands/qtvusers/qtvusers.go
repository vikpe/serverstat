package qtvusers

import (
	"errors"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'},
}

func ParseResponse(responseBody []byte, err error) ([]qclient.Client, error) {
	if err != nil {
		return []qclient.Client{}, err
	} else {
		return ParseResponseBody(responseBody)
	}
}

func ParseResponseBody(responseBody []byte) ([]qclient.Client, error) {
	// example response body: 12 "djevulsk" "serp" "player" "rst" "twitch.tv/vikpe"
	fullText := string(responseBody)
	const QuoteChar = "\""

	if !strings.Contains(fullText, QuoteChar) {
		return []qclient.Client{}, errors.New("invalid response body")
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
