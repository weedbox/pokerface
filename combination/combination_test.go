package combination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPossibleCombinations_Standard(t *testing.T) {

	cards := []string{"S2", "H3", "D4", "C5", "C7", "DT", "DK"}

	combinations := GetAllPossibleCombinations(cards, 5)

	for _, c := range combinations {
		assert.Equal(t, len(c), 5)
	}
}

func TestGetAllPossibleCombinations_NotEnoughSources(t *testing.T) {

	cards := []string{"S2", "H3", "D4", "C5"}

	combinations := GetAllPossibleCombinations(cards, 5)

	assert.Equal(t, len(combinations), 1)

	for _, c := range combinations {
		assert.Equal(t, len(c), 4)
	}
}
