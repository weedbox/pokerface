package pot

import (
	"sort"
)

type LevelList struct {
	contributors  map[int]int64
	foldedPlayers map[int]bool
	levels        []*Level
}

func NewLevelList() *LevelList {
	return &LevelList{
		contributors:  make(map[int]int64),
		foldedPlayers: make(map[int]bool),
		levels:        make([]*Level, 0),
	}
}

func (ll *LevelList) Count() int {
	return len(ll.levels)
}

func (ll *LevelList) GetLevels() []*Level {
	return ll.levels
}

func (ll *LevelList) AssertLevel(level int64) *Level {

	for _, pot := range ll.levels {

		if pot.Level == level {
			return pot
		}
	}

	// Not found pot so create a new one
	l := &Level{
		Level:        level,
		Wager:        0,
		Total:        0,
		Contributors: make([]int, 0),
	}

	ll.levels = append(ll.levels, l)

	return l
}

func (ll *LevelList) AddContributor(wager int64, contributorIdx int, fold bool) {

	ll.contributors[contributorIdx] = wager
	if fold {
		ll.foldedPlayers[contributorIdx] = true
	}

	ll.AssertLevel(wager)

	sort.Slice(ll.levels, func(i, j int) bool {
		return ll.levels[i].Level < ll.levels[j].Level
	})

	// Add contributers to each levels
	for _, pot := range ll.levels {

		// Reset contributor list
		pot.Contributors = make([]int, 0)
		for idx, wager := range ll.contributors {
			if pot.Level <= wager {
				pot.Contributors = append(pot.Contributors, idx)
			}
		}

	}

	// Calculate total wagers for each levels
	prevLevel := int64(0)
	for _, l := range ll.levels {
		l.Wager = l.Level - prevLevel
		l.Total = int64(len(l.Contributors)) * l.Wager
		prevLevel = l.Level
	}
}

func (ll *LevelList) GetPots() []*Pot {

	// Preparing new levels its contributors doesn't contain foldded players
	origPots := make([]*Pot, 0)
	for _, l := range ll.levels {

		p := &Pot{
			Level:        l.Level,
			Wager:        l.Wager,
			Total:        l.Total,
			Contributors: make(map[int]int64),
			Levels:       make([]*Level, 0),
		}

		p.Levels = append(p.Levels, l)

		// Arrage contributors
		for _, cIdx := range l.Contributors {

			// Remove folded player
			if _, ok := ll.foldedPlayers[cIdx]; ok {
				continue
			}

			p.Contributors[cIdx] = l.Wager
		}

		origPots = append(origPots, p)
	}

	// Merge pots which contains the same contributors
	pots := make([]*Pot, 0)
	var prev *Pot = nil
	for i, p := range origPots {

		if i == 0 {
			pots = append(pots, p)
			prev = p
			continue
		}

		// Comparing contributors
		if len(prev.Contributors) != len(p.Contributors) {
			pots = append(pots, p)
			prev = p
			continue
		}

		prev.Level = p.Level
		prev.Wager += p.Wager
		prev.Total += p.Total
		prev.Levels = append(prev.Levels, p.Levels...)

		for pIdx, wager := range p.Contributors {
			prev.Contributors[pIdx] += wager
		}
	}

	// Put foldded players back to pots
	for pIdx, _ := range ll.foldedPlayers {

		wager := ll.contributors[pIdx]
		if wager == 0 {
			continue
		}

		for _, p := range pots {
			p.Contributors[pIdx] = wager

			if wager < p.Level {
				break
			}
		}
	}

	return pots
}
