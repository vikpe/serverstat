package qwfwd

import (
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

const Name = "qwfwd"
const VersionPrefix = Name

type Qwfwd struct {
	Address     string             `json:"address"`
	ClientNames []string           `json:"client_names"`
	Settings    qsettings.Settings `json:"settings"`
	Geo         geo.Info           `json:"geo"`
}
