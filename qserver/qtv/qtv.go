package qtv

import (
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
)

const Name = "qtv"
const VersionPrefix = Name

type Qtv struct {
	Address        string
	SpectatorNames []qstring.QuakeString
	Settings       qsettings.Settings
	Geo            geo.Info
}

type QtvExport struct {
	Type string
	Qtv
}

func Export(qtv Qtv) QtvExport {
	return QtvExport{
		Type: Name,
		Qtv:  qtv,
	}
}
