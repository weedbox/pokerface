package main

import (
	"sort"

	"github.com/cfsghost/pokerface/pot"
	"github.com/cfsghost/pokerface/settlement"
)

type RankInfo struct {
	Player *PlayerState
	Power  *PowerState
}

func (g *game) GetAlivePlayers() []*PlayerState {

	players := make([]*PlayerState, 0)
	for _, p := range g.gs.Players {

		// Find the player who did not fold
		if !p.Fold {
			players = append(players, &p)
		}
	}

	return players
}

func (g *game) CalculatePlayersRanking() []*RankInfo {

	players := g.GetAlivePlayers()

	// Calculate power for all players
	ranks := make([]*RankInfo, 0)
	for _, p := range players {
		powerState := g.CalculatePlayerPower(p)

		ranks = append(ranks, &RankInfo{
			Player: p,
			Power:  powerState,
		})
	}

	// Sort by power score
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].Power.Score > ranks[j].Power.Score
	})

	return ranks
}

func (g *game) CalculateGameResults() error {

	r := &settlement.Result{
		Players: make([]*settlement.PlayerResult, 0),
		Pots:    make([]*settlement.PotResult, 0),
	}

	// Initializing player results
	for _, p := range g.gs.Players {
		pr := &settlement.PlayerResult{
			Idx:     p.Idx,
			Final:   p.Bankroll,
			Changed: 0,
		}

		r.Players = append(r.Players, pr)
	}

	// Initializing pot results
	for _, pot := range g.gs.Status.Pots {
		r.AddPot(pot.Total)
	}

	// Update winner and loser results for each pot
	for potIdx, pot := range g.gs.Status.Pots {

		potRank := g.CalculatePotRank(pot)
		winners := potRank.GetWinners()

		// Calculate chips for multiple winners of this pot
		chips := pot.Total / int64(len(winners))

		//TODO: Solve problem that chips of pot is indivisible by winners

		for _, wIdx := range winners {
			r.Withdraw(potIdx, wIdx, chips)
		}

		// Update loser results (should be negtive)
		losers := potRank.GetLoser()
		for _, lIdx := range losers {
			r.Withdraw(potIdx, lIdx, -pot.Wager)
		}
	}

	return nil
}

func (g *game) CalculatePotRank(p *pot.Pot) *pot.PotRank {

	pr := pot.NewPotRank()
	for _, c := range p.Contributers {
		ps := g.Player(c).State()
		if ps.Fold {
			pr.AddContributer(0, c)
			continue
		}

		pr.AddContributer(g.Player(c).State().Combination.Power, c)
	}

	pr.Calculate()

	return pr
}
