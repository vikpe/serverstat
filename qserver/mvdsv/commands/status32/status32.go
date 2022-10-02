package status32

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', ' '},
}

func ParseResponse(address string, responseBody []byte) (qtvstream.QtvStream, error) {
	// example response body
	// ����nqtv 1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 0
	reader := csv.NewReader(strings.NewReader(string(responseBody)))
	reader.Comma = ' '
	reader.FieldsPerRecord = 4

	record, err := reader.Read()
	if err != nil {
		return qtvstream.New(), errors.New("unable to parse response")
	}

	const (
		IndexTitle       = 1
		IndexUrl         = 2
		IndexClientCount = 3
	)

	streamTitle := record[IndexTitle]
	streamURL := record[IndexUrl]

	if streamURL == "" {
		// invalid configuration (not reporting qtv ip as they should)
		// example response body: 1 "DuelMania FRANCE Qtv (1)" "" 0
		// try to parse stream number from title
		streamNumber, err := StreamNumberFromTitle(streamTitle)

		if err != nil {
			return qtvstream.New(), errors.New("invalid QTV configuration")
		}

		hostname, _, _ := net.SplitHostPort(address)
		streamURL = fmt.Sprintf("%d@%s:%d", streamNumber, hostname, qtv.PortNumber)
	}

	urlParts := strings.Split(streamURL, "@")

	stream := qtvstream.QtvStream{
		Title:          streamTitle,
		Url:            streamURL,
		ID:             qutil.StringToInt(urlParts[0]),
		Address:        urlParts[1],
		SpectatorNames: make([]string, 0),
		SpectatorCount: qutil.StringToInt(record[IndexClientCount]),
	}
	return stream, nil
}

func StreamNumberFromTitle(title string) (int, error) {
	indexOpenBrace := strings.LastIndex(title, "(")
	indexCloseBrace := strings.LastIndex(title, ")")
	err := errors.New("unable to parse stream number from title")

	if -1 == indexOpenBrace || -1 == indexCloseBrace {
		return 0, err
	}

	numberAsString := title[indexOpenBrace+1 : indexCloseBrace]

	if 0 == len(numberAsString) {
		return 0, err
	}

	return qutil.StringToInt(numberAsString), nil
}
