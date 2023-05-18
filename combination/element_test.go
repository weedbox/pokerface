package combination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetElementsByRank(t *testing.T) {

	cards := GetCardStates([]string{"S2", "H5", "D9", "CT", "CK"})
	elements := GetElementsByRank(cards)

	assert.Equal(t, len(elements), 5)
}

func TestGetElementsByRank_Pair(t *testing.T) {

	cards := GetCardStates([]string{"S2", "H2", "D9", "CT", "CK"})
	elements := GetElementsByRank(cards)

	assert.Equal(t, len(elements), 4)

	count := 5
	for _, ele := range elements {
		assert.LessOrEqual(t, ele.Count, count)
		count = ele.Count
	}
}

func TestGetElementsByRank_TwoPair(t *testing.T) {

	cards := GetCardStates([]string{"S2", "H2", "DT", "CT", "CK"})
	elements := GetElementsByRank(cards)

	assert.Equal(t, len(elements), 3)

	count := 5
	for _, ele := range elements {
		assert.LessOrEqual(t, ele.Count, count)
		count = ele.Count
	}
}

func TestGetElementsByRank_ThreeOfAKind(t *testing.T) {

	cards := GetCardStates([]string{"S2", "H2", "D2", "C3", "CK"})
	elements := GetElementsByRank(cards)

	assert.Equal(t, len(elements), 3)

	count := 5
	for _, ele := range elements {
		assert.LessOrEqual(t, ele.Count, count)
		count = ele.Count
	}
}

func TestGetElementsByRank_FullHouse(t *testing.T) {

	cards := GetCardStates([]string{"S2", "H2", "D2", "C3", "S3"})
	elements := GetElementsByRank(cards)

	assert.Equal(t, len(elements), 2)

	count := 5
	for _, ele := range elements {
		assert.LessOrEqual(t, ele.Count, count)
		count = ele.Count
	}
}

func TestGetElementsByRank_FourOfAKind(t *testing.T) {

	cards := GetCardStates([]string{"S2", "H2", "D2", "C2", "S3"})
	elements := GetElementsByRank(cards)

	assert.Equal(t, len(elements), 2)

	count := 5
	for _, ele := range elements {
		assert.LessOrEqual(t, ele.Count, count)
		count = ele.Count
	}
}
