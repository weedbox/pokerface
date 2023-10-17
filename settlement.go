package pokerface

import (
	"github.com/weedbox/pokerface/combination"
	"github.com/weedbox/pokerface/settlement"
)

type RankInfo struct {
	Player *PlayerState
	Power  *combination.PowerState
}

/*
func (g *game) GetAlivePlayers() []*PlayerState {

		players := make([]*PlayerState, 0)
		for _, p := range g.gs.Players {

			// Find the player who did not fold
			if !p.Fold {
				players = append(players, p)
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
*/
func (g *game) CalculateGameResults() error {

	r := settlement.NewResult()

	// Initializing pot results
	for _, pot := range g.gs.Status.Pots {
		r.AddPot(pot.Total, pot.Levels)
	}

	// Initializing player scores
	for _, p := range g.gs.Players {

		r.AddPlayer(p.Idx, p.Bankroll)

		// No score if player fold already
		if p.Fold {
			r.UpdateScore(p.Idx, 0)
			continue
		}

		r.UpdateScore(p.Idx, p.Combination.Power)
	}

	r.Calculate()

	// Update state
	g.gs.Result = r

	return nil
}
