package main

import "sort"

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

func (g *game) CalculatePowerRanking() []*RankInfo {

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
