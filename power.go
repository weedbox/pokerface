package pokerface

import (
	"sort"

	"github.com/cfsghost/pokerface/combination"
)

func (g *game) CalculatePlayerPower(p *PlayerState) *combination.PowerState {

	// calculate power with player state
	powers := g.GetAllPowersByPlayer(p)

	// The first combination is the best result
	return powers[0]
}

func (g *game) UpdateCombinationOfAllPlayers() error {

	for _, p := range g.gs.Players {
		ps := g.CalculatePlayerPower(p)

		p.Combination.Type = combination.CombinationSymbol[ps.Combination]

		for _, c := range ps.Cards {
			p.Combination.Cards = append(p.Combination.Cards, c.ToString())
		}

		p.Combination.Power = int(ps.Score)
	}

	return nil
}

func (g *game) GetAllPowersByPlayer(p *PlayerState) []*combination.PowerState {

	powers := make([]*combination.PowerState, 0)

	// Calcuate power for all combinations
	combinations := g.GetAllPossibileCombinations(p, g.gs.Meta.RequiredHoleCardsCount)
	for _, c := range combinations {
		ps := g.CalculateCombinationPower(c)
		powers = append(powers, ps)
	}

	// sorting combination by power score
	sort.Slice(powers, func(i, j int) bool {
		return powers[i].Score > powers[j].Score
	})

	return powers
}

func (g *game) CalculateCombinationPower(cards []string) *combination.PowerState {
	return combination.CalculatePower(g.gs.Meta.CombinationPowers, cards)
}

func (g *game) GetAllPossibileCombinations(p *PlayerState, holeCardsCount int) [][]string {

	combinations := make([][]string, 0)

	if holeCardsCount == 0 {
		allCards := make([]string, 0)
		allCards = append(allCards, p.HoleCards...)
		allCards = append(allCards, g.gs.Status.Board...)
		return combination.GetAllPossibleCombinations(allCards, 5)
	}

	holeCardCombinations := combination.GetAllPossibleCombinations(p.HoleCards, holeCardsCount)
	boardCardCombinations := combination.GetAllPossibleCombinations(g.gs.Status.Board, 5-holeCardsCount)

	for _, cards := range holeCardCombinations {
		allCards := make([]string, 0)
		allCards = append(allCards, cards...)

		for _, bCards := range boardCardCombinations {
			allCards = append(allCards, bCards...)
		}

		combinations = append(combinations, allCards)
	}

	return combinations
}
