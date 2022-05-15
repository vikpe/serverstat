package qversion

import "strings"

type Type struct {
	Name          string
	VersionPrefix string
}

var (
	TypeFte         = Type{Name: "fte", VersionPrefix: "fte"}
	TypeMvdsv       = Type{Name: "mvdsv", VersionPrefix: "mvdsv"}
	TypeQwfwd       = Type{Name: "qwfwd", VersionPrefix: "qwfwd"}
	TypeQtv         = Type{Name: "qtv", VersionPrefix: "qtv"}
	TypeFortressOne = Type{Name: "fortress_one", VersionPrefix: "fo svn"}
	TypeUnknown     = Type{Name: "unknown", VersionPrefix: ""}
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

func IsType(version string, serverType Type) bool {
	return strings.HasPrefix(
		strings.ToLower(version),
		strings.ToLower(serverType.VersionPrefix),
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
