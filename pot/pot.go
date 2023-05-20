package pot

import "sort"

type Pot struct {
	Wager        int64 `json:"wager"`
	Total        int64 `json:"total"`
	Contributers []int `json:"contributers"`
}

type PotList struct {
	pots []*Pot
}

func NewPotList() *PotList {
	return &PotList{
		pots: make([]*Pot, 0),
	}
}

func (p *PotList) Count() int {
	return len(p.pots)
}

func (p *PotList) GetPots() []*Pot {
	return p.pots
}

func (p *PotList) AddContributer(wager int64, contributerIdx int) {

	for _, pot := range p.pots {

		if pot.Wager == wager {
			pot.Total += wager
			pot.Contributers = append(pot.Contributers, contributerIdx)
			return
		}
	}

	// Not found pot so create a new one
	pot := &Pot{
		Wager:        wager,
		Total:        wager,
		Contributers: make([]int, 0),
	}

	pot.Contributers = append(pot.Contributers, contributerIdx)

	p.pots = append(p.pots, pot)

	// Sort by wager
	sort.Slice(p.pots, func(i, j int) bool {
		return p.pots[i].Wager < p.pots[j].Wager
	})
}
