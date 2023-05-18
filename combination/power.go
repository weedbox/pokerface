package combination

import (
	"math"
	"sort"
)

type PowerState struct {
	Combination Combination
	Score       uint64
	Cards       []*Card
	Elements    []*Element
}

type CombinationPower struct {
	Combination Combination
	Score       uint64
}

func CalculatePower(pr PowerRankings, cardSymbols []string) *PowerState {

	// Transform card strings to internal structure
	cards := GetCardStates(cardSymbols)

	// Sorting based on card rank
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank > cards[j].Rank
	})

	ps := &PowerState{
		Cards:       cards,
		Combination: CombinationHighCard,
		Elements:    GetElementsByRank(cards),
	}

	// Flush
	if isFlush(cards) {
		ps.Combination = CombinationFlush
	}

	// Straight
	if isStraight(cards) {
		if ps.Combination == CombinationFlush {
			ps.Combination = CombinationStraightFlush
		} else {
			ps.Combination = CombinationStraight
		}
	}

	// Other combinations
	if isFourOfAKind(ps.Elements) {
		ps.Combination = CombinationFourOfAKind
	} else if isFullHouse(ps.Elements) {
		ps.Combination = CombinationFullHouse
	} else if isThreeOfAKind(ps.Elements) {
		ps.Combination = CombinationThreeOfAKind
	} else if isTwoPair(ps.Elements) {
		ps.Combination = CombinationTwoPair
	} else if isPair(ps.Elements) {
		ps.Combination = CombinationPair
	}

	powerBaseline := CalculatePowerLevels(pr, ps)
	score := CalculatePowerScore(ps)
	ps.Score = score + powerBaseline

	//fmt.Printf("raw_score=%d, level_power=%d\n", score, powerBaseline)

	return ps
}

func CalculatePowerLevels(pr PowerRankings, ps *PowerState) uint64 {

	powerLevel := uint64(0)
	for _, c := range pr {

		if ps.Combination == c {
			return powerLevel
		}

		powerLevel += CombinationLevel[c]
	}

	return 0
}

func CalculatePowerScore(ps *PowerState) uint64 {

	score := uint64(0)

	switch ps.Combination {
	case CombinationStraight:
		fallthrough
	case CombinationStraightFlush:

		totalPoint := 0
		maxRank := 0
		for _, e := range ps.Elements {
			if maxRank < e.Rank {
				maxRank = e.Rank
			}

			totalPoint += e.Rank
		}

		// A, 2, 3, 4, 5
		if maxRank == 14 && totalPoint == 28 {
			score = 0
		} else {
			// >= 2, 3, 4, 5, 6
			score = uint64(maxRank) - 5
		}

	default:
		for i, e := range ps.Elements {
			level := len(ps.Elements) - i - 1
			based := math.Pow(13, float64(level))
			s := uint64((float64)(e.Rank-2) * based) // calibration for Ace(14) and alignment to 0, so reduce
			score += s
			//			fmt.Printf("rank=%d, level=%d, based=%f, score=%d, total_score=%d\n", e.Rank, level, based, s, score)
		}
	}

	return score
}

func isFlush(cards []*Card) bool {

	if len(cards) == 0 {
		return false
	}

	suit := cards[0].Suit

	for _, c := range cards {
		if suit != c.Suit {
			return false
		}
	}

	return true
}

func isStraight(cards []*Card) bool {

	if len(cards) != 5 {
		return false
	}

	// No chance to be straight if highest rank is less than 5
	if cards[0].Rank < 5 {
		return false
	}

	restOfCards := cards

	// The highest rank is Ace(14) that could be two types of straight
	if cards[0].Rank == 14 && cards[1].Rank == 5 {
		// assume that lowest rank of straight
		restOfCards = cards[1:5]
	}

	// Check each rank
	cur := restOfCards[0].Rank
	for _, c := range restOfCards {
		if c.Rank != cur {
			return false
		}

		cur--
	}

	return true
}

func isFourOfAKind(elements []*Element) bool {

	for _, ele := range elements {
		if ele.Count == 4 {
			return true
		}
	}

	return false
}

func isFullHouse(elements []*Element) bool {

	hasThree := false
	hasTwo := false
	for _, ele := range elements {
		if ele.Count == 3 {
			hasThree = true
		}
		if ele.Count == 2 {
			hasTwo = true
		}
	}
	return hasThree && hasTwo
}

func isThreeOfAKind(elements []*Element) bool {

	for _, ele := range elements {
		if ele.Count == 3 {
			return true
		}
	}
	return false
}

func isTwoPair(elements []*Element) bool {

	pairCount := 0
	for _, ele := range elements {
		if ele.Count == 2 {
			pairCount++
		}
	}

	if pairCount == 2 {
		return true
	}

	return false
}

func isPair(elements []*Element) bool {

	pairCount := 0
	for _, ele := range elements {
		if ele.Count == 2 {
			pairCount++
		}
	}

	if pairCount == 1 {
		return true
	}

	return false
}
