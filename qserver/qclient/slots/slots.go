package slots

type Slots struct {
	Used  int
	Total int
	Free  int
}

func New(total int, used int) Slots {
	return Slots{
		Used:  used,
		Total: total,
		Free:  total - used,
	}
}
