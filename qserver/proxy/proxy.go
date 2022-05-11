package proxy

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
)

type Proxy struct {
	Address     string
	ClientNames []qstring.QuakeString
	Settings    qsettings.Settings
}

func Parse(genericServer qserver.GenericServer) Proxy {
	clientNames := make([]qstring.QuakeString, 0)

	for _, client := range genericServer.Clients {
		clientNames = append(clientNames, client.Name)
	}

	return Proxy{
		Address:     genericServer.Address,
		ClientNames: clientNames,
		Settings:    genericServer.Settings,
	}
}
