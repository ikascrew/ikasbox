package api

type Status struct {
	success bool
}

func (s Status) IsSuccess() bool {
	return s.success
}
