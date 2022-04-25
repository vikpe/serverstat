package serverstat

type Player struct {
	Name    string
	NameInt []int
	Team    string
	TeamInt []int
	Skin    string
	Colors  [2]int
	Frags   int
	Ping    int
	Time    int
	IsBot   bool
}

type Client struct {
	Player
	IsSpec bool
}

type Spectator struct {
	Name    string
	NameInt []int
	IsBot   bool
}

type QtvStream struct {
	Id            int
	Title         string
	Url           string
	NumSpectators int
}

func newQtvStream() QtvStream {
	return QtvStream{
		Id:            0,
		Title:         "",
		Url:           "",
		NumSpectators: 0,
	}
}

type QtvServer struct {
	Title         string
	Address       string
	NumSpectators int
}

type QuakeServer struct {
	Title         string
	Address       string
	QtvStream     QtvStream
	Map           string
	NumPlayers    int
	MaxPlayers    int
	NumSpectators int
	MaxSpectators int
	Players       []Player
	Spectators    []Spectator
	Settings      map[string]string
}

func newQuakeServer() QuakeServer {
	return QuakeServer{
		Title:      "",
		Address:    "",
		Settings:   map[string]string{},
		Players:    make([]Player, 0),
		Spectators: make([]Spectator, 0),
		QtvStream:  newQtvStream(),
	}
}
