package laststats

// todo
import (
	"bytes"
	"errors"
	"github.com/goccy/go-json"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/udpclient"
)

var FrameDelimiter = []byte{0xff, 0xff, 0xff, 0xff, 'n'}

func GetCommand(limit int) udpclient.Command {
	return udpclient.Command{
		RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' ', byte(limit), 0x0a},
		ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' '},
	}
}

func ParseResponseBody(responseBody []byte, err error) ([]Entry, error) {
	// example: "����nlaststats 2\n[json..]\n"
	if err != nil {
		return []Entry{}, err
	}

	// remove frame delimiters
	responseBody = bytes.ReplaceAll(responseBody, FrameDelimiter, []byte{})

	// validate body
	jsonIndexBegin := bytes.Index(responseBody, []byte("["))
	jsonIndexEnd := bytes.LastIndex(responseBody, []byte("]"))

	if -1 == jsonIndexBegin || -1 == jsonIndexEnd || jsonIndexBegin > jsonIndexEnd {
		return []Entry{}, errors.New("malformed json")
	}

	// empty body
	if jsonIndexBegin+1 == jsonIndexEnd {
		return []Entry{}, nil
	}

	// non empty body
	jsonBody := responseBody[jsonIndexBegin : jsonIndexEnd+1]
	cleanJsonBody := []byte(qutil.StripControlCharacters(string(jsonBody)))

	var entries []Entry
	err = json.Unmarshal(cleanJsonBody, &entries)

	if err != nil {
		return []Entry{}, err
	}

	if len(entries) > 0 && 0 == len(entries[0].Date) {
		return []Entry{}, errors.New("invalid fields, date is missing")
	}

	return entries, nil
}
