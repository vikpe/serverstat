package qtv

import (
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
)

const Name = "qtv"
const VersionPrefix = Name

type Qtv struct {
	Address        string
	Type           string
	SpectatorNames []qstring.QuakeString
	Settings       qsettings.Settings
}
