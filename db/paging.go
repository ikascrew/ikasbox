package db

type Paging struct {
	Current int `json:"current"`
	Count   int `json:"count"`
	Limit   int `json:"limit"`
}

func (p *Paging) SetCount(cnt int) {
	p.Count = cnt
}
