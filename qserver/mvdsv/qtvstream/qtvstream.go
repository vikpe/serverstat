package qtvstream

type QtvStream struct {
	Title          string   `json:"title"`
	Url            string   `json:"url"`
	ID             int      `json:"id"`
	Address        string   `json:"address"`
	SpectatorNames []string `json:"spectator_names"`
	SpectatorCount int      `json:"spectator_count"`
}

func New() QtvStream {
	return QtvStream{
		SpectatorNames: make([]string, 0),
	}
}
