package qtv

import (
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

const Name = "qtv"
const VersionPrefix = Name

type Qtv struct {
	Address        string
	SpectatorNames []string
	Settings       qsettings.Settings
	Geo            geo.Info
}
