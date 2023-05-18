package combination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePower_CombinationPowerStandard(t *testing.T) {

	cardSets := [][]string{
		// HighCard
		[]string{"S2", "H3", "D4", "C5", "C7"},
		[]string{"S2", "H3", "D4", "C5", "CT"},
		[]string{"SA", "HK", "DQ", "CJ", "C9"},

		// Pair
		[]string{"S2", "H2", "D4", "C5", "C7"},
		[]string{"SA", "HA", "D4", "C5", "C7"},
		[]string{"SA", "HA", "DJ", "CQ", "CK"},

		// TwoPair
		[]string{"S2", "H2", "D4", "C3", "C3"},
		[]string{"S2", "H2", "D4", "C4", "C7"},
		[]string{"SA", "HA", "D4", "CK", "SK"},
		[]string{"SA", "HA", "DQ", "CK", "SK"},

		// ThreeOfAKind
		[]string{"S2", "H2", "D2", "C4", "C7"},
		[]string{"SK", "HK", "DK", "C4", "C7"},
		[]string{"SA", "HA", "DA", "CK", "CQ"},

		// Straight
		[]string{"SA", "H2", "D3", "C4", "C5"},
		[]string{"S2", "H3", "D4", "C5", "C6"},
		[]string{"ST", "HJ", "DQ", "CK", "CA"},

		// Flush
		[]string{"C2", "C3", "C4", "C5", "C7"},
		[]string{"CT", "C3", "C4", "C5", "C6"},
		[]string{"CA", "C3", "C4", "C5", "C6"},
		[]string{"CA", "C3", "C4", "C5", "C7"},
		[]string{"CA", "C3", "C4", "C5", "CT"},
		[]string{"CA", "C3", "C4", "C5", "CK"},
		[]string{"CA", "C9", "CJ", "CQ", "CK"},

		// FullHouse
		[]string{"S2", "H2", "D2", "C3", "S3"},
		[]string{"S3", "H3", "D3", "C4", "S4"},
		[]string{"S3", "H3", "D3", "C5", "S5"},
		[]string{"S3", "H3", "D3", "CA", "SA"},
		[]string{"SK", "HK", "DK", "CA", "SA"},
		[]string{"SA", "HA", "DA", "CT", "ST"},
		[]string{"SA", "HA", "DA", "CK", "SK"},

		// FourOfAKind
		[]string{"S2", "H2", "D2", "C2", "C4"},
		[]string{"S2", "H2", "D2", "C2", "C5"},
		[]string{"S2", "H2", "D2", "C2", "CT"},
		[]string{"S2", "H2", "D2", "C2", "CA"},
		[]string{"S3", "H3", "D3", "C3", "CA"},
		[]string{"SK", "HK", "DK", "CK", "CA"},
		[]string{"SA", "HA", "DA", "CA", "CK"},

		// StraightFlush
		[]string{"S2", "S3", "S4", "S5", "SA"},
		[]string{"S2", "S3", "S4", "S5", "S6"},
		[]string{"S3", "S4", "S5", "S6", "S7"},
		[]string{"ST", "SJ", "SQ", "SK", "SA"},
	}

	prevScore := uint64(0)
	for _, cardSymbols := range cardSets {
		ps := CalculatePower(CombinationPowerStandard, cardSymbols)

		//t.Log(cardSymbols, CombinationSymbol[ps.Combination], ps.Score)

		// Should be greater than previous score
		assert.Greater(t, ps.Score, prevScore)

		prevScore = ps.Score
	}
}

func TestCalculatePower_CombinationPowerShortDeck(t *testing.T) {

	cardSets := [][]string{
		// HighCard
		[]string{"S2", "H3", "D4", "C5", "C7"},
		[]string{"S2", "H3", "D4", "C5", "CT"},
		[]string{"SA", "HK", "DQ", "CJ", "C9"},

		// Pair
		[]string{"S2", "H2", "D4", "C5", "C7"},
		[]string{"SA", "HA", "D4", "C5", "C7"},
		[]string{"SA", "HA", "DJ", "CQ", "CK"},

		// TwoPair
		[]string{"S2", "H2", "D4", "C3", "C3"},
		[]string{"S2", "H2", "D4", "C4", "C7"},
		[]string{"SA", "HA", "D4", "CK", "SK"},
		[]string{"SA", "HA", "DQ", "CK", "SK"},

		// ThreeOfAKind
		[]string{"S2", "H2", "D2", "C4", "C7"},
		[]string{"SK", "HK", "DK", "C4", "C7"},
		[]string{"SA", "HA", "DA", "CK", "CQ"},

		// Straight
		[]string{"SA", "H2", "D3", "C4", "C5"},
		[]string{"S2", "H3", "D4", "C5", "C6"},
		[]string{"ST", "HJ", "DQ", "CK", "CA"},

		// FullHouse
		[]string{"S2", "H2", "D2", "C3", "S3"},
		[]string{"S3", "H3", "D3", "C4", "S4"},
		[]string{"S3", "H3", "D3", "C5", "S5"},
		[]string{"S3", "H3", "D3", "CA", "SA"},
		[]string{"SK", "HK", "DK", "CA", "SA"},
		[]string{"SA", "HA", "DA", "CT", "ST"},
		[]string{"SA", "HA", "DA", "CK", "SK"},

		// Flush
		[]string{"C2", "C3", "C4", "C5", "C7"},
		[]string{"CT", "C3", "C4", "C5", "C6"},
		[]string{"CA", "C3", "C4", "C5", "C6"},
		[]string{"CA", "C3", "C4", "C5", "C7"},
		[]string{"CA", "C3", "C4", "C5", "CT"},
		[]string{"CA", "C3", "C4", "C5", "CK"},
		[]string{"CA", "C9", "CJ", "CQ", "CK"},

		// FourOfAKind
		[]string{"S2", "H2", "D2", "C2", "C4"},
		[]string{"S2", "H2", "D2", "C2", "C5"},
		[]string{"S2", "H2", "D2", "C2", "CT"},
		[]string{"S2", "H2", "D2", "C2", "CA"},
		[]string{"S3", "H3", "D3", "C3", "CA"},
		[]string{"SK", "HK", "DK", "CK", "CA"},
		[]string{"SA", "HA", "DA", "CA", "CK"},

		// StraightFlush
		[]string{"S2", "S3", "S4", "S5", "SA"},
		[]string{"S2", "S3", "S4", "S5", "S6"},
		[]string{"S3", "S4", "S5", "S6", "S7"},
		[]string{"ST", "SJ", "SQ", "SK", "SA"},
	}

	prevScore := uint64(0)
	for _, cardSymbols := range cardSets {
		ps := CalculatePower(CombinationPowerShortDeck, cardSymbols)

		//t.Log(cardSymbols, CombinationSymbol[ps.Combination], ps.Score)

		// Should be greater than previous score
		assert.Greater(t, ps.Score, prevScore)

		prevScore = ps.Score
	}
}

func TestCalculatePower_Partial(t *testing.T) {

	// HighCard
	ps := CalculatePower(CombinationPowerStandard, []string{"S2", "HA"})
	assert.Equal(t, ps.Combination, CombinationHighCard)
	assert.Equal(t, len(ps.Elements), 2)

	// Pair
	ps = CalculatePower(CombinationPowerStandard, []string{"S2", "H2"})
	assert.Equal(t, ps.Combination, CombinationPair)
	assert.Equal(t, len(ps.Elements), 1)

	// HighCard
	ps = CalculatePower(CombinationPowerStandard, []string{"S2", "H3", "D4"})
	assert.Equal(t, ps.Combination, CombinationHighCard)
	assert.Equal(t, len(ps.Elements), 3)

	// Pair
	ps = CalculatePower(CombinationPowerStandard, []string{"S2", "H2", "D4"})
	assert.Equal(t, ps.Combination, CombinationPair)
	assert.Equal(t, len(ps.Elements), 2)

	// ThreeOfAKind
	ps = CalculatePower(CombinationPowerStandard, []string{"S2", "H2", "D2"})
	assert.Equal(t, ps.Combination, CombinationThreeOfAKind)
	assert.Equal(t, len(ps.Elements), 1)

	// TwoPair
	ps = CalculatePower(CombinationPowerStandard, []string{"S2", "H2", "D3", "C3"})
	assert.Equal(t, ps.Combination, CombinationTwoPair)
	assert.Equal(t, len(ps.Elements), 2)

	// FourOfAKind
	ps = CalculatePower(CombinationPowerStandard, []string{"S2", "H2", "D2", "C2"})
	assert.Equal(t, ps.Combination, CombinationFourOfAKind)
	assert.Equal(t, len(ps.Elements), 1)
}

func TestCalculatePower_HighCard(t *testing.T) {

	cardSets := [][]string{
		[]string{"S2", "H3", "D4", "C5", "C7"},
		[]string{"S2", "H3", "D4", "C5", "C8"},
		[]string{"S2", "H3", "D4", "C5", "C9"},
		[]string{"S2", "H3", "D4", "C5", "CT"},
		[]string{"S2", "H3", "D4", "C5", "CJ"},
		[]string{"S2", "H3", "D4", "C5", "CQ"},
		[]string{"S2", "H3", "D4", "C5", "CK"},
		[]string{"S2", "H3", "D4", "C6", "CA"},
		[]string{"S2", "H3", "D4", "C7", "CA"},
		[]string{"S2", "H3", "D4", "C8", "CA"},
		[]string{"S2", "H3", "D4", "C9", "CA"},
		[]string{"S2", "H3", "D4", "CT", "CA"},
		[]string{"S2", "H3", "D4", "CJ", "CA"},
		[]string{"S2", "H3", "D4", "CQ", "CA"},
		[]string{"S2", "H3", "D4", "CK", "CA"},
		[]string{"S2", "H3", "D5", "CK", "CA"},
		[]string{"S2", "H3", "D6", "CK", "CA"},
		[]string{"S2", "H3", "D7", "CK", "CA"},
		[]string{"S2", "H3", "D8", "CK", "CA"},
		[]string{"S2", "H3", "D9", "CK", "CA"},
		[]string{"S2", "H3", "DT", "CK", "CA"},
		[]string{"S2", "H4", "DT", "CK", "CA"},
		[]string{"S2", "H5", "DT", "CK", "CA"},
		[]string{"S2", "H6", "DT", "CK", "CA"},
		[]string{"S2", "H7", "DT", "CK", "CA"},
		[]string{"S2", "H8", "DT", "CK", "CA"},
		[]string{"S3", "H8", "DT", "CK", "CA"},
		[]string{"S4", "H8", "DT", "CK", "CA"},
		[]string{"S5", "H8", "DT", "CK", "CA"},
	}

	prevScore := uint64(0)
	for _, cardSymbols := range cardSets {
		ps := CalculatePower(CombinationPowerStandard, cardSymbols)

		assert.Equal(t, ps.Combination, CombinationHighCard)
		assert.Equal(t, len(ps.Elements), 5)

		// Should be greater than previous score
		assert.Greater(t, ps.Score, prevScore)
	}
}

func TestCalculatePower_Pair(t *testing.T) {

	cardSets := [][]string{
		[]string{"S2", "H2", "D4", "C5", "C7"},
		[]string{"S2", "H2", "D4", "C5", "C8"},
		[]string{"S2", "H2", "D4", "C5", "C9"},
		[]string{"S2", "H2", "D4", "C5", "CT"},
		[]string{"S2", "H2", "D4", "C5", "CA"},
		[]string{"S2", "H2", "D4", "C9", "CA"},
		[]string{"S2", "H2", "D5", "C9", "CA"},
		[]string{"S2", "H2", "D6", "C9", "CA"},
	}

	prevScore := uint64(0)
	for _, cardSymbols := range cardSets {
		ps := CalculatePower(CombinationPowerStandard, cardSymbols)

		assert.Equal(t, ps.Combination, CombinationPair)
		assert.Equal(t, len(ps.Elements), 4)

		// Should be greater than previous score
		assert.Greater(t, ps.Score, prevScore)
	}
}

func TestCalculatePower_Straight(t *testing.T) {

	cardSets := [][]string{
		[]string{"SA", "H2", "D3", "C4", "C5"},
		[]string{"S5", "H6", "D7", "C8", "C9"},
		[]string{"ST", "HJ", "DQ", "CK", "CA"},
	}

	prevScore := uint64(0)
	for _, cardSymbols := range cardSets {
		ps := CalculatePower(CombinationPowerStandard, cardSymbols)

		assert.Equal(t, ps.Combination, CombinationStraight)
		assert.Equal(t, len(ps.Elements), 5)

		// Should be greater than previous score
		assert.Greater(t, ps.Score, prevScore)
	}
}

func TestCalculatePower_Straight_Invalid(t *testing.T) {

	// J, Q, K, A, 2 is high card rather than straight
	cardSymbols := []string{"SJ", "HQ", "DK", "CA", "C2"}

	ps := CalculatePower(CombinationPowerStandard, cardSymbols)

	assert.Equal(t, ps.Combination, CombinationHighCard)
	assert.Equal(t, len(ps.Elements), 5)
}
