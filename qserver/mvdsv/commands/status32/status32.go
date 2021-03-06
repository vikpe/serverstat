package status32

import (
	"encoding/csv"
	"errors"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', ' '},
}

func ParseResponse(responseBody []byte, err error) (qtvstream.QtvStream, error) {
	if err != nil {
		return qtvstream.New(), err
	}

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

	if record[IndexUrl] == "" {
		// invalid configuration (not reporting qtv ip as they should)
		return qtvstream.New(), errors.New("invalid QTV configuration")
	}

	urlParts := strings.Split(record[IndexUrl], "@")

	stream := qtvstream.QtvStream{
		Title:          record[IndexTitle],
		Url:            record[IndexUrl],
		ID:             qutil.StringToInt(urlParts[0]),
		Address:        urlParts[1],
		SpectatorNames: make([]string, 0),
		SpectatorCount: qutil.StringToInt(record[IndexClientCount]),
	}
	return stream, nil
}
