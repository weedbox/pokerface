package pot

import (
	"sort"
)

type PotList struct {
	contributors map[int]int64
	pots         []*Pot
}

func NewPotList() *PotList {
	return &PotList{
		contributors: make(map[int]int64),
		pots:         make([]*Pot, 0),
	}
}

func (p *PotList) Count() int {
	return len(p.pots)
}

func (p *PotList) GetPots() []*Pot {
	return p.pots
}

func (p *PotList) AssertPot(level int64) *Pot {

	for _, pot := range p.pots {

		if pot.Level == level {
			return pot
		}
	}

	// Not found pot so create a new one
	pot := &Pot{
		Level:        level,
		Wager:        0,
		Total:        0,
		Contributors: make([]int, 0),
	}

	p.pots = append(p.pots, pot)

	return pot
}

func (p *PotList) AddContributor(wager int64, contributorIdx int) {

	p.contributors[contributorIdx] = wager

	p.AssertPot(wager)

	sort.Slice(p.pots, func(i, j int) bool {
		return p.pots[i].Level < p.pots[j].Level
	})

	// Add contributers to pots
	for _, pot := range p.pots {

		// Reset contributor list
		pot.Contributors = make([]int, 0)
		for idx, wager := range p.contributors {
			if pot.Level <= wager {
				pot.Contributors = append(pot.Contributors, idx)
			}
		}

	}

	// Calculate total wagers of pots
	prevLevel := int64(0)
	for _, pot := range p.pots {
		pot.Wager = pot.Level - prevLevel
		pot.Total = int64(len(pot.Contributors)) * pot.Wager
		prevLevel = pot.Level
	}
}
