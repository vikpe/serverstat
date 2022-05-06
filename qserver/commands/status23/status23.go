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

func Send(udpClient udpclient.UdpClient, address string) (qserver.GenericServer, error) {
	// request
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}
	server, err := ParseResponse(udpClient.Request(address, statusPacket, expectedHeader))

	if err != nil {
		server.Address = address
	}

	return server, nil
}

func ParseResponse(responseBody []byte, err error) (qserver.GenericServer, error) {
	if err != nil {
		return qserver.GenericServer{}, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))
	scanner.Scan()

	settingsString := scanner.Text()

	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	server := qserver.GenericServer{
		Version:  qversion.New(qsettings.New(settingsString)["*version"]),
		Clients:  qclient.NewFromStrings(clientStrings),
		Settings: qsettings.New(settingsString),
	}
	return server, nil
}
