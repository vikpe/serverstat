package slots

type Slots struct {
	Used  int `json:"used"`
	Total int `json:"total"`
	Free  int `json:"free"`
}

func New(total int, used int) Slots {
	return Slots{
		Used:  used,
		Total: total,
		Free:  total - used,
	}
}
