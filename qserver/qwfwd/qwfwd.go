package qwfwd

import (
	"encoding/json"

	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
)

const Name = "qwfwd"
const VersionPrefix = Name

type Qwfwd struct {
	Address     string
	ClientNames []qstring.QuakeString
	Settings    qsettings.Settings
}

func (server Qwfwd) MarshalJSON() ([]byte, error) {
	type qwfwdJson struct {
		Address     string
		Type        string
		ClientNames []qstring.QuakeString
		Settings    qsettings.Settings
	}

	return json.Marshal(qwfwdJson{
		Address:     server.Address,
		Type:        Name,
		ClientNames: server.ClientNames,
		Settings:    server.Settings,
	})
}
