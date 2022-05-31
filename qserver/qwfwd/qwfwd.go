package qwfwd

import (
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

const Name = "qwfwd"
const VersionPrefix = Name

type Qwfwd struct {
	Address     string
	ClientNames []string
	Settings    qsettings.Settings
	Geo         geo.Info
}
