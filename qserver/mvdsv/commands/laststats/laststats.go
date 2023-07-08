package laststats

// todo
import (
	"bytes"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/vikpe/serverstat/qutil"
	"github.com/vikpe/udpclient"
)

var FrameDelimiter = []byte{0xff, 0xff, 0xff, 0xff, 'n'}

func GetCommand(limit int) udpclient.Command {
	return udpclient.Command{
		RequestPacket:  append([]byte{0xff, 0xff, 0xff, 0xff}, []byte(fmt.Sprintf("laststats %d\n", limit))...),
		ResponseHeader: append([]byte{0xff, 0xff, 0xff, 0xff}, []byte("nlaststats ")...),
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
	jsonBody := cleanJson(responseBody[jsonIndexBegin : jsonIndexEnd+1])
	var entries []Entry
	err = json.Unmarshal(jsonBody, &entries)

	if err != nil {
		return []Entry{}, err
	}

	if len(entries) > 0 && 0 == len(entries[0].Date) {
		return []Entry{}, errors.New("invalid fields, date is missing")
	}

	// convert player/team names from unicode to ascii
	for entryIndex, entry := range entries {
		for playerIndex, player := range entry.Players {
			entries[entryIndex].Players[playerIndex].Name = qutil.UnicodeToAscii(player.Name)
			entries[entryIndex].Players[playerIndex].Team = qutil.UnicodeToAscii(player.Team)
		}

		for teamIndex, teamName := range entry.Teams {
			entries[entryIndex].Teams[teamIndex] = qutil.UnicodeToAscii(teamName)
		}
	}

	return entries, nil
}

func cleanJson(value []byte) []byte {
	result := []byte(qutil.StripControlCharacters(string(value)))
	result = bytes.ReplaceAll(result, []byte(",]"), []byte("]"))
	return bytes.ReplaceAll(result, []byte("[,"), []byte("["))
}
