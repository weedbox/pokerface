package settlement

import (
	"github.com/weedbox/pokerface/pot"
)

type Result struct {
	Players []*PlayerResult `json:"players"`
	Pots    []*PotResult    `json:"pots"`
}

type PlayerResult struct {
	Idx     int   `json:"idx"`
	Final   int64 `json:"final"`
	Changed int64 `json:"changed"`
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

func (r *Result) AddPot(total int64, levels []*pot.Level) {

	pr := &PotResult{
		level:   NewPotLevel(),
		Total:   total,
		Winners: make([]*Winner, 0),
	}

	for _, l := range levels {
		pr.level.AddLevel(l.Level, l.Wager, l.Total, l.Contributors)
	}

	r.Pots = append(r.Pots, pr)
}

func (r *Result) UpdateScore(playerIdx int, score int) {

	for _, p := range r.Pots {
		for _, l := range p.level.levels {
			l.UpdateScore(playerIdx, score)
		}
	}
}

func (r *Result) Update(potIdx int, playerIdx int, wager int64, withdraw int64) {

	pot := r.Pots[potIdx]

	// Update winners information
	if withdraw > 0 {
		pot.UpdateWinner(playerIdx, withdraw+wager)
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

func (r *Result) CalculateWinnerRewards(potIdx int, l *LevelInfo) {

	// Calculate contributer ranks of this pot by score
	l.rank.Calculate()

	// Calculate chips for multiple winners of this pot
	winners := l.rank.GetWinners()

	// Calculate rewards
	based := l.Total / int64(len(winners))
	remainder := l.Total % int64(len(winners))

	for i, wIdx := range winners {

		reward := based

		if int64(i) < remainder {
			reward += 1
		}

		r.Update(potIdx, wIdx, l.Wager, reward-l.Wager)
	}
}

func (r *Result) CalculateLoserResults(potIdx int, l *LevelInfo) {

	losers := l.rank.GetLoser()

	for _, lIdx := range losers {
		// withdraw should be negtive
		r.Update(potIdx, lIdx, l.Wager, -l.Wager)
	}
}

func (r *Result) CalculatePot(potIdx int, p *PotResult) {

	for _, l := range p.level.levels {

		// Calculate chips for multiple winners of this pot
		r.CalculateWinnerRewards(potIdx, l)

		// Update loser results
		r.CalculateLoserResults(potIdx, l)
	}
}

func (r *Result) Calculate() {

	for potIdx, pot := range r.Pots {
		r.CalculatePot(potIdx, pot)
	}
}
