package main

type Combination int32

const (
	CombinationHighCard Combination = iota
	CombinationPair
	CombinationTwoPair
	CombinationThreeOfAKind
	CombinationStraight
	CombinationFlush
	CombinationFullHouse
	CombinationFourOfAKind
	CombinationStraightFlush
)
