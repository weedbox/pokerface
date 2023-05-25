package pokerface

import (
	"github.com/cfsghost/pokerface/combination"
	"github.com/cfsghost/pokerface/settlement"
)

type RankInfo struct {
	Player *PlayerState
	Power  *combination.PowerState
}

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

/*
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

	// Initializing player results
	for _, p := range g.gs.Players {
		r.AddPlayer(p.Idx, p.StackSize)
	}

	// Initializing pot results
	for _, pot := range g.gs.Status.Pots {
		r.AddPot(pot.Total)
	}

	// Add contributers to each pot
	for potIdx, pot := range g.gs.Status.Pots {

		for _, c := range pot.Contributors {
			player := g.Player(c).State()

			// No score if player fold already
			if player.Fold {
				r.AddContributor(potIdx, c, 0)
				continue
			}

			r.AddContributor(potIdx, c, player.Combination.Power)
		}
	}

	r.Calculate()

	// Update state
	g.gs.Result = r

	return nil
}
