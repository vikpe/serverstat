package quakeserver

import "github.com/vikpe/qw-serverstat/qtvstream"

type Player struct {
	Name    string
	NameRaw []byte
	Team    string
	TeamRaw []byte
	Skin    string
	Colors  [2]uint8
	Frags   uint16
	Ping    uint16
	Time    uint8
	IsBot   bool
}

type Client struct {
	Player
	IsSpec bool
}

type Spectator struct {
	Name    string
	NameRaw []byte
	IsBot   bool
}

type QuakeServer struct {
	Title         string
	Address       string
	QtvStream     qtvstream.QtvStream
	Map           string
	NumPlayers    uint8
	MaxPlayers    uint8
	NumSpectators uint8
	MaxSpectators uint8
	Players       []Player
	Spectators    []Spectator
	Settings      map[string]string
}

func New() QuakeServer {
	return QuakeServer{
		Title:      "",
		Address:    "",
		Settings:   map[string]string{},
		Players:    make([]Player, 0),
		Spectators: make([]Spectator, 0),
		QtvStream:  qtvstream.New(),
	}
}
