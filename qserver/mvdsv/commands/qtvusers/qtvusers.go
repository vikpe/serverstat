package qtvusers

import (
	"errors"
	"strings"

	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 'q', 't', 'v', 'u', 's', 'e', 'r', 's', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'},
}

func ParseResponse(responseBody []byte, err error) ([]string, error) {
	if err != nil {
		return make([]string, 0), err
	}

	// example response body: 12 "djevulsk" "serp" "player" "rst" "twitch.tv/vikpe"
	fullText := string(responseBody)
	const QuoteChar = "\""

	if !strings.Contains(fullText, QuoteChar) {
		return make([]string, 0), errors.New("invalid response body")
	}

	indexFirstQuote := strings.Index(fullText, QuoteChar)
	indexLastQuote := strings.LastIndex(fullText, QuoteChar)
	namesText := fullText[indexFirstQuote+1 : indexLastQuote]
	spectatorQuakeNames := strings.Split(namesText, "\" \"")

	spectatorPlainNames := make([]string, 0)
	for _, quakeName := range spectatorQuakeNames {
		spectatorPlainNames = append(spectatorPlainNames, qstring.New(quakeName).ToPlainString())
	}

	return spectatorPlainNames, nil
}
