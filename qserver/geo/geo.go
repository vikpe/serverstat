package geo

type Location struct {
	CC          string     `json:"cc"`
	Country     string     `json:"country"`
	Region      string     `json:"region"`
	City        string     `json:"city"`
	Coordinates [2]float32 `json:"coordinates"`
}
