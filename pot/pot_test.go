package pot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPotList(t *testing.T) {

	contributers := []int64{
		1000,
		1000,
		1000,
		2000,
		2000,
		3000,
	}

	list := NewPotList()

	for idx, wager := range contributers {
		list.AddContributer(wager, idx)
	}

	assert.Equal(t, 3, list.Count())

	pots := list.GetPots()

	assert.Equal(t, int64(1000), pots[0].Wager)
	assert.Equal(t, int64(2000), pots[1].Wager)
	assert.Equal(t, int64(3000), pots[2].Wager)

	prevPotWager := int64(0)
	for _, p := range pots {
		assert.Greater(t, p.Wager, prevPotWager)
	}

	assert.Equal(t, 3, len(pots[0].Contributers))
	assert.Equal(t, 2, len(pots[1].Contributers))
	assert.Equal(t, 1, len(pots[2].Contributers))

	assert.Equal(t, int64(3000), pots[0].Total)
	assert.Equal(t, int64(4000), pots[1].Total)
	assert.Equal(t, int64(3000), pots[2].Total)
}
