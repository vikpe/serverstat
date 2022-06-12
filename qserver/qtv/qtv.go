package qtv

import (
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

const Name = "qtv"
const VersionPrefix = Name

type Qtv struct {
	Address        string             `json:"address"`
	SpectatorNames []string           `json:"spectator_names"`
	Settings       qsettings.Settings `json:"settings"`
	Geo            geo.Info           `json:"geo"`
}
