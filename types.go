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

type QtvServer struct {
	Title         string
	Address       string
	Numspectators int
	Spectators    []string
}

type QuakeServer struct {
	Title         string
	Address       string
	QtvAddress    string
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
		QtvAddress: "",
	}
}
