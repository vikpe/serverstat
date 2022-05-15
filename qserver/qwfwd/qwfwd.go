package qwfwd

import (
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
