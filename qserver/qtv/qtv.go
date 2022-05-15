package qtv

import (
	"encoding/json"

	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
)

const Name = "qtv"
const VersionPrefix = Name

type Qtv struct {
	Address        string
	SpectatorNames []qstring.QuakeString
	Settings       qsettings.Settings
}

func (server Qtv) MarshalJSON() ([]byte, error) {
	type qtvJson struct {
		Address        string
		Type           string
		SpectatorNames []qstring.QuakeString
		Settings       qsettings.Settings
	}

	return json.Marshal(qtvJson{
		Address:        server.Address,
		Type:           Name,
		SpectatorNames: server.SpectatorNames,
		Settings:       server.Settings,
	})
}
