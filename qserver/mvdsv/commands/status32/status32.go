package status32

import (
	"encoding/csv"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/udpclient"
)

func SendTo(address string) (qtvstream.QtvStream, error) {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v'}
	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return qtvstream.QtvStream{}, err
	}

	responseBody := response[5:]
	return ParseResponseBody(responseBody, err)
}

func ParseResponseBody(responseBody []byte, err error) (qtvstream.QtvStream, error) {
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

	stream := qtvstream.QtvStream{
		Title:      record[IndexTitle],
		Url:        record[IndexAddress],
		NumClients: uint8(qutil.StringToInt(record[IndexClientCount])),
	}
	return stream, nil
}
