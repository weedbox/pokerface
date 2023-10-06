package combination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPossibleCombinations_Standard(t *testing.T) {

	cards := []string{"S2", "H3", "D4", "C5", "C7", "DT", "DK"}

	combinations := GetPossibleCombinations(cards, 5)

	for _, c := range combinations {
		assert.Equal(t, 5, len(c))
	}
}

func TestGetPossibleCombinations_NotEnoughSources(t *testing.T) {

	cards := []string{"S2", "H3", "D4", "C5"}

	combinations := GetPossibleCombinations(cards, 5)

	assert.Equal(t, len(combinations), 1)

	for _, c := range combinations {
		assert.Equal(t, 4, len(c))
	}
}

func TestGetAllPossibleCombinations_Standard(t *testing.T) {

	board := []string{"D6", "ST", "H9", "S6"}
	hole := []string{"H7", "CQ"}

	combinations := GetAllPossibleCombinations(board, hole, 2)

	for _, c := range combinations {
		assert.Equal(t, 5, len(c))
	}
}

func TestGetAllPossibleCombinations_4_HoleCards(t *testing.T) {

	board := []string{"D6", "ST", "H9", "S6"}
	hole := []string{"H7", "CQ", "CK", "DT"}

	combinations := GetAllPossibleCombinations(board, hole, 2)

	for _, c := range combinations {
		assert.Equal(t, 5, len(c))
	}
}
