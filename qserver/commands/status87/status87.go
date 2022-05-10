package status87

import (
	"bufio"
	"strings"

	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '8', '7', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'},
}

func ParseResponse(responseBody []byte, err error) (map[string]string, []qclient.Client, error) {
	if err != nil {
		return map[string]string{}, []qclient.Client{}, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))
	scanner.Scan()

	settingsString := scanner.Text()

	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	settings := qsettings.ParseString(settingsString)
	clients := qclient.NewFromStrings(clientStrings)

	return settings, clients, nil
}
