package qversion

import (
	"strings"

	"github.com/vikpe/serverstat/qserver/fortressone"
	"github.com/vikpe/serverstat/qserver/fte"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

type Version string

func New(value string) Version {
	return Version(value)
}

func (v Version) IsMvdsv() bool {
	return v.hasPrefix(mvdsv.VersionPrefix)
}

func (v Version) IsFte() bool {
	return v.hasPrefix(fte.VersionPrefix)
}

func (v Version) IsQwfwd() bool {
	return v.hasPrefix(qwfwd.VersionPrefix)
}

func (v Version) IsQtv() bool {
	return v.hasPrefix(qtv.VersionPrefix)
}

func (v Version) IsFortressOne() bool {
	return v.hasPrefix(fortressone.VersionPrefix)
}

func (v Version) hasPrefix(prefix string) bool {
	return strings.HasPrefix(strings.ToLower(string(v)), strings.ToLower(prefix))
}

func (v Version) GetType() string {
	if v.IsMvdsv() {
		return mvdsv.Name
	} else if v.IsQwfwd() {
		return qwfwd.Name
	} else if v.IsQtv() {
		return qtv.Name
	} else if v.IsFte() {
		return fte.Name
	} else if v.IsFortressOne() {
		return fortressone.Name
	} else {
		return "unknown"
	}
}
