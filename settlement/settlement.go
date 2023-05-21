package settlement

type Result struct {
	Players []*PlayerResult `json:"players"`
	Pots    []*PotResult    `json:"pots"`
}

type PlayerResult struct {
	Idx     int   `json:"idx"`
	Final   int64 `json:"final"`
	Changed int64 `json:"changed"`
}

type PotResult struct {
	Chips   int64     `json:"chips"`
	Winners []*Winner `json:"winners"`
}

type Winner struct {
	Idx   int   `json:"idx"`
	Chips int64 `json:"chips"`
}

func (r *Result) AddPot(total int64) {

	pr := &PotResult{
		Chips:   total,
		Winners: make([]*Winner, 0),
	}

	r.Pots = append(r.Pots, pr)
}

func (r *Result) Withdraw(potIdx int, playerIdx int, chips int64) {

	// Update pot result
	pot := r.Pots[potIdx]

	if chips >= 0 {

		// Add winner to pot
		w := &Winner{
			Idx:   playerIdx,
			Chips: chips,
		}

		pot.Winners = append(pot.Winners, w)
	}

	for _, p := range r.Players {
		if p.Idx == playerIdx {
			p.Final += chips
			p.Changed += chips
			return
		}
	}
}
