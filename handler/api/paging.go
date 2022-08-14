package api

type Paging struct {
	Current int `"json":"current"`
	Count   int `"json":"count"`
	Limit   int `"json":"limit"`
	MaxPage int `"json":"maxPage"`
}

func (p *Paging) SetCount(cnt int) {
	p.Count = cnt
	p.MaxPage = 0

	if p.Limit <= 0 {
		if p.Count != 0 {
			p.MaxPage = 1
		}
		return
	}

	max := cnt / p.Limit
	if (cnt % p.Limit) != 0 {
		max++
	}
	p.MaxPage = max
}
