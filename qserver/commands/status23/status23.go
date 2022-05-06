package status23

import (
	"bufio"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'},
}

func ParseResponse(responseBody []byte, err error) (map[string]string, []qclient.Client, error) {
	if err != nil {
		return nil, nil, err
	} else {
		settings, clients := ParseResponseBody(responseBody)
		return settings, clients, nil
	}
}

func ParseResponseBody(responseBody []byte) (map[string]string, []qclient.Client) {
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))
	scanner.Scan()

	settingsString := scanner.Text()

	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	settings := qsettings.New(settingsString)
	clients := qclient.NewFromStrings(clientStrings)

	return settings, clients
}
