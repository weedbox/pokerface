package main

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

func (g *game) CalculatePower(p *PlayerState) *PowerState {

	//TODO: calculate power with player state

	return &PowerState{}
}
