package main

import "sort"

var CombinationPowerStandard = []Combination{
	CombinationHighCard,
	CombinationPair,
	CombinationTwoPair,
	CombinationThreeOfAKind,
	CombinationStraight,
	CombinationFlush,
	CombinationFullHouse,
	CombinationFourOfAKind,
	CombinationStraightFlush,
}

var CombinationPowerShortDeck = []Combination{
	CombinationHighCard,
	CombinationPair,
	CombinationTwoPair,
	CombinationThreeOfAKind,
	CombinationStraight,
	CombinationFullHouse,
	CombinationFlush,
	CombinationFourOfAKind,
	CombinationStraightFlush,
}

type PowerState struct {
	Combination Combination
	Score       uint64
	Cards       []string
}

func (g *game) CalculatePlayerPower(p *PlayerState) *PowerState {

	// calculate power with player state
	powers := g.GetAllCombinationsByPlayer(p)

	// The first combination is the best result
	return powers[0]
}

func (g *game) GetAllCombinationsByPlayer(p *PlayerState) []*PowerState {

	powers := make([]*PowerState, 0)
	combinations := make([][]string, 0)

	if g.gs.Meta.RequiredHoleCardsCount == 0 {
		//TODO:
		combinations = append(combinations, []string{})
	} else {
		//TODO: must pick cards from hole cards
		combinations = append(combinations, []string{})
	}

	// Calcuate power for all combinations
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

func (g *game) CalculateCombinationPower(cards []string) *PowerState {

	//TODO: Calculate score

	return &PowerState{
		Cards: cards,
	}
}
