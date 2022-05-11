package qtv

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
)

type Qtv struct {
	Address        string
	SpectatorNames []qstring.QuakeString
	Settings       qsettings.Settings
}

func Parse(genericServer qserver.GenericServer) Qtv {
	spectatorNames := make([]qstring.QuakeString, 0)

	for _, client := range genericServer.Clients {
		spectatorNames = append(spectatorNames, client.Name)
	}

	return Qtv{
		Address:        genericServer.Address,
		SpectatorNames: spectatorNames,
		Settings:       genericServer.Settings,
	}
}
