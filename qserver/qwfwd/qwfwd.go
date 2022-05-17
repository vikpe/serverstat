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

type QwfwdExport struct {
	Type string
	Qwfwd
}

func Export(qwfwd Qwfwd) QwfwdExport {
	return QwfwdExport{
		Type:  Name,
		Qwfwd: qwfwd,
	}
}
