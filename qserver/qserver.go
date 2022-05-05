package qserver

import (
	"bufio"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/udpclient"
)

type GenericServer struct {
	Version   qversion.Version `json:"-"`
	Address   string
	Clients   []qclient.Client
	Settings  map[string]string
	ExtraInfo extraInfo `json:"-"`
}

type extraInfo struct {
	QtvStream mvdsv.QtvStream
}

func New(address string) (GenericServer, error) {
	// request
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}

	udpClient := udpclient.New()
	response, err := udpClient.Request(address, statusPacket, expectedHeader)

	if err != nil {
		return GenericServer{}, err
	}

	// response
	responseBody := response[len(expectedHeader):]
	settingsString, clientStrings := parseStatusResponseBody(responseBody)

	// result
	settings := qsettings.New(settingsString)
	clients := qclient.NewFromStrings(clientStrings)
	version := qversion.New(settings["*version"])
	extraInfo := getExtraInfo(version, address)

	server := GenericServer{
		Version:   version,
		Address:   address,
		Clients:   clients,
		Settings:  settings,
		ExtraInfo: extraInfo,
	}

	return server, nil
}

func parseStatusResponseBody(responseBody []byte) (string, []string) {
	scanner := bufio.NewScanner(strings.NewReader(string(responseBody)))

	scanner.Scan()
	settingsString := scanner.Text()

	var clientStrings []string
	for scanner.Scan() {
		clientStrings = append(clientStrings, scanner.Text())
	}

	return settingsString, clientStrings
}

func getExtraInfo(version qversion.Version, address string) extraInfo {
	extraInfo := extraInfo{}

	if version.IsMvdsv() {
		qtvStream, err := mvdsv.GetQtvStreamInfo(address)

		if err != nil {
			extraInfo.QtvStream = qtvStream
		}
	}

	return extraInfo
}
