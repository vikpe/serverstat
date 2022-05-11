package status32

import (
	"encoding/csv"
	"errors"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '3', '2', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', ' '},
}

func ParseResponse(responseBody []byte, err error) (qtvstream.QtvStream, error) {
	if err != nil {
		return qtvstream.QtvStream{}, err
	}

	// example repsonse body
	// ����nqtv 1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 0
	reader := csv.NewReader(strings.NewReader(string(responseBody)))
	reader.Comma = ' '
	reader.FieldsPerRecord = 4

	record, err := reader.Read()
	if err != nil {
		return qtvstream.QtvStream{}, errors.New("unable to parse response")
	}

	const (
		IndexTitle       = 1
		IndexAddress     = 2
		IndexClientCount = 3
	)

	if record[IndexAddress] == "" {
		// invalid configuration (not reporting qtv ip as they should)
		return qtvstream.QtvStream{}, errors.New("invalid QTV configuration")
	}

	stream := qtvstream.QtvStream{
		Title:          record[IndexTitle],
		Url:            record[IndexAddress],
		NumSpectators:  uint8(qutil.StringToInt(record[IndexClientCount])),
		SpectatorNames: []qstring.QuakeString{},
	}
	return stream, nil
}
