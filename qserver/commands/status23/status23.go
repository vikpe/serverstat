package status23

import (
	"bufio"
	"strings"

	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/udpclient"
)

func SendTo(address string) (qserver.GenericServer, error) {
	// request
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}

	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return qserver.GenericServer{}, err
	}

	// response
	responseBody := response[len(expectedHeader):]
	server := ParseResponseBody(responseBody)
	server.Address = address

	return server, nil
}

func ParseResponseBody(responseBody []byte) qserver.GenericServer {
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))

	scanner.Scan()
	settingsString := scanner.Text()

	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	return qserver.GenericServer{
		Version:  qversion.New(qsettings.New(settingsString)["*version"]),
		Clients:  qclient.NewFromStrings(clientStrings),
		Settings: qsettings.New(settingsString),
	}
}
