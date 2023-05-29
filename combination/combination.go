package combination

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

var CombinationSymbol = map[Combination]string{
	CombinationHighCard:      "HighCard",
	CombinationPair:          "Pair",
	CombinationTwoPair:       "TwoPair",
	CombinationThreeOfAKind:  "ThreeOfAKind",
	CombinationStraight:      "Straight",
	CombinationFlush:         "Flush",
	CombinationFullHouse:     "FullHouse",
	CombinationFourOfAKind:   "FourOfAKind",
	CombinationStraightFlush: "StraightFlush",
}

// Power score of combination:
// HighCard			: 13 * 12 * 11 * 10 * 9 ~= 13^5				= 371,293
// Pair				: 13(pair) * 12(pair) * 11 * 10 ~= 13^4		= 28,561
// TwoPair			: 13(pair) * 12(pair) * 11 ~= 13^3			= 2,197
// ThreeOfAKind		: 13(toak) * 12 * 11 ~= 13^3				= 2,197
// Straight			: 10(5~A) ~= 13^1							= 13
// Flush			: 13 * 12 * 11 * 10 * 9 ~= 13^5				= 371,293
// FullHouse		: 13(toak) * 12(pair) ~= 13^2				= 169
// FourOfAKind		: 13(foak) * 12 ~= 13^2						= 169
// StraightFlush	: 10(5~A) ~= 13^1							= 13

var CombinationLevel = map[Combination]uint64{
	CombinationHighCard:      371293,
	CombinationPair:          28561,
	CombinationTwoPair:       2197,
	CombinationThreeOfAKind:  2197,
	CombinationStraight:      13,
	CombinationFlush:         371293,
	CombinationFullHouse:     169,
	CombinationFourOfAKind:   169,
	CombinationStraightFlush: 13,
}

// Combination ranking table for different game rules
type PowerRankings []Combination

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

func gospersHack(k int, n int) []int {

	result := make([]int, 0)

	cur := (1 << k) - 1
	limit := 1 << n
	for cur < limit {
		// do something
		//fmt.Printf("%0*b\n", n, cur)
		result = append(result, cur)

		lb := cur & -cur
		r := cur + lb
		cur = (((r ^ cur) >> 2) / lb) | r
	}

	return result
}

func binaryOnesPositions(value int, n int) []int {

	var positions []int
	for i := 0; i < n; i++ {
		if (value>>i)&1 == 1 {
			positions = append(positions, i)
		}
	}
	return positions
}

func GetPossibleCombinations(cards []string, n int) [][]string {

	combinations := make([][]string, 0)

	total := len(cards)
	if total <= n {
		combinations = append(combinations, cards)
		return combinations
	}

	posBins := gospersHack(n, total)

	for _, v := range posBins {
		positions := binaryOnesPositions(v, total)
		combination := make([]string, 0)
		for _, p := range positions {
			combination = append(combination, cards[p])
		}

		combinations = append(combinations, combination)
	}

	return combinations
}

func GetAllPossibleCombinations(boardCards []string, holeCards []string, holeCardsCount int) [][]string {

	combinations := make([][]string, 0)

	if holeCardsCount == 0 {
		allCards := make([]string, 0)
		allCards = append(allCards, holeCards...)
		allCards = append(allCards, boardCards...)
		return GetPossibleCombinations(allCards, 5)
	}

	holeCardCombinations := GetPossibleCombinations(holeCards, holeCardsCount)
	boardCardCombinations := GetPossibleCombinations(boardCards, 5-holeCardsCount)

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
