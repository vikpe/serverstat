package laststats

// todo
import (
	"bytes"
	"errors"

	"github.com/goccy/go-json"
	"github.com/vikpe/udpclient"
)

var frameDelimiter = []byte{0xff, 0xff, 0xff, 0xff, 'n'}

func GetCommand(limit int) udpclient.Command {
	return udpclient.Command{
		RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' ', byte(limit), 0x0a},
		ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' '},
	}
}

func ParseResponseBody(responseBody []byte, err error) ([]Entry, error) {
	// example: "2\n[json..]\n"
	if err != nil {
		return []Entry{}, err
	}

	// remove frame delimiters
	responseBody = bytes.ReplaceAll(responseBody, frameDelimiter, []byte{})

	// validate body
	jsonError := errors.New("invalid json in response body")
	jsonIndexBegin := bytes.Index(responseBody, []byte("["))
	jsonIndexEnd := bytes.LastIndex(responseBody, []byte("]"))

	if -1 == jsonIndexBegin || -1 == jsonIndexEnd || jsonIndexBegin > jsonIndexEnd {
		return []Entry{}, jsonError
	}

	// parse body
	if jsonIndexBegin+1 == jsonIndexEnd {
		return []Entry{}, nil
	}

	jsonBody := responseBody[jsonIndexBegin : jsonIndexEnd+1]
	var stats []Entry

	err = json.Unmarshal(jsonBody, &stats)
	if err != nil {
		return []Entry{}, jsonError
	}

	if len(stats) > 0 && 0 == len(stats[0].Date) {
		return []Entry{}, jsonError
	}

	return stats, nil
}
