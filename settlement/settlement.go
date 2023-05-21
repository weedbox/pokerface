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
	rank PotRank

	Chips   int64     `json:"chips"`
	Winners []*Winner `json:"winners"`
}

type Winner struct {
	Idx   int   `json:"idx"`
	Chips int64 `json:"chips"`
}

func NewResult() *Result {
	return &Result{
		Players: make([]*PlayerResult, 0),
		Pots:    make([]*PotResult, 0),
	}
}

func (r *Result) AddPlayer(playerIdx int, bankroll int64) {

	pr := &PlayerResult{
		Idx:     playerIdx,
		Final:   bankroll,
		Changed: 0,
	}

	r.Players = append(r.Players, pr)
}

func (r *Result) AddPot(total int64) {

	pr := &PotResult{
		Chips:   total,
		Winners: make([]*Winner, 0),
	}

	r.Pots = append(r.Pots, pr)
}

func (r *Result) AddContributer(potIdx int, playerIdx int, score int) {

	// Take pot
	pot := r.Pots[potIdx]

	// Add a new contributr
	pot.rank.AddContributer(score, playerIdx)
}

func (r *Result) Update(potIdx int, playerIdx int, withdraw int64) {

	pot := r.Pots[potIdx]

	// Add winner to pot
	if withdraw >= 0 {

		w := &Winner{
			Idx:   playerIdx,
			Chips: withdraw,
		}

		pot.Winners = append(pot.Winners, w)
	}

	// Update player results
	for _, p := range r.Players {
		if p.Idx == playerIdx {
			p.Final += withdraw
			p.Changed += withdraw
			return
		}
	}
}

func (r *Result) CalculateWagerOfPot(total int64, contributerCount int) int64 {
	return total / int64(contributerCount)
}

func (r *Result) CalculateWinnerRewards(potIdx int, wager int64, total int64, winners []int) {

	based := total / int64(len(winners))
	remainder := total % int64(len(winners))

	for i, wIdx := range winners {

		reward := based

		if int64(i) < remainder {
			reward += 1
		}

		r.Update(potIdx, wIdx, reward-wager)
	}
}

func (r *Result) Calculate() {

	for potIdx, pot := range r.Pots {

		// Calculate contributions of players
		wager := r.CalculateWagerOfPot(pot.Chips, pot.rank.ContributerCount())

		// Calculate contributer ranks of this pot by score
		pot.rank.Calculate()

		// Calculate chips for multiple winners of this pot
		winners := pot.rank.GetWinners()
		r.CalculateWinnerRewards(potIdx, wager, pot.Chips, winners)

		// Update loser results (should be negtive)
		losers := pot.rank.GetLoser()
		for _, lIdx := range losers {
			r.Update(potIdx, lIdx, -wager)
		}
	}
}
