package qversion

import "strings"

type Type string

const (
	TypeFte     Type = "fte"
	TypeMvdsv   Type = "mvdsv"
	TypeProxy   Type = "qwfwd"
	TypeQtv     Type = "qtv"
	TypeUnknown Type = "unknown"
)

type Version string

func New(value string) Version {
	return Version(value)
}

func (v Version) IsMvdsv() bool {
	return v.IsType(TypeMvdsv)
}

func (v Version) IsFte() bool {
	return v.IsType(TypeFte)
}

func (v Version) IsProxy() bool {
	return v.IsType(TypeProxy)
}

func (v Version) IsQtv() bool {
	return v.IsType(TypeQtv)
}

func (v Version) IsGameServer() bool {
	return v.IsMvdsv() || v.IsFte()
}

func (v Version) IsType(t Type) bool {
	return IsType(string(v), t)
}

func (v Version) GetType() Type {
	return GetType(string(v))
}

func IsMvdsv(version string) bool {
	return IsType(version, TypeMvdsv)
}

func IsFte(version string) bool {
	return IsType(version, TypeFte)
}

func IsProxy(version string) bool {
	return IsType(version, TypeProxy)
}

func IsQtv(version string) bool {
	return IsType(version, TypeQtv)
}

func IsGameServer(version string) bool {
	return IsMvdsv(version) || IsFte(version)
}

func IsType(version string, serverType string) bool {
	return strings.Contains(
		strings.ToLower(version),
		strings.ToLower(serverType),
	)
}

func GetType(v string) Type {
	if IsProxy(v) {
		return TypeProxy
	} else if IsMvdsv(v) {
		return TypeMvdsv
	} else if IsFte(v) {
		return TypeFte
	} else if IsQtv(v) {
		return TypeQtv
	} else {
		return TypeUnknown
	}
}
