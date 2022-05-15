package qversion

import "strings"

type Type string

const (
	TypeFte         Type = "fte"
	TypeMvdsv       Type = "mvdsv"
	TypeQwfwd       Type = "qwfwd"
	TypeQtv         Type = "qtv"
	TypeFortressOne Type = "fo svn"
	TypeUnknown     Type = "unknown"
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

func (v Version) IsQwfwd() bool {
	return v.IsType(TypeQwfwd)
}

func (v Version) IsQtv() bool {
	return v.IsType(TypeQtv)
}

func (v Version) IsGameServer() bool {
	return IsGameServer(string(v))
}

func (v Version) IsFortressOne() bool {
	return v.IsType(TypeFortressOne)
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

func IsQwfwd(version string) bool {
	return IsType(version, TypeQwfwd)
}

func IsQtv(version string) bool {
	return IsType(version, TypeQtv)
}

func IsFortressOne(version string) bool {
	return IsType(version, TypeFortressOne)
}

func IsGameServer(version string) bool {
	return IsMvdsv(version) || IsFte(version)
}

func IsType(version string, serverType Type) bool {
	return strings.Contains(
		strings.ToLower(version),
		strings.ToLower(string(serverType)),
	)
}

func GetType(v string) Type {
	if IsMvdsv(v) {
		return TypeMvdsv
	} else if IsQwfwd(v) {
		return TypeQwfwd
	} else if IsQtv(v) {
		return TypeQtv
	} else if IsFte(v) {
		return TypeFte
	} else if IsFortressOne(v) {
		return TypeFortressOne
	} else {
		return TypeUnknown
	}
}
