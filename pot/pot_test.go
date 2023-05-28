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
		list.AddContributor(wager, idx)
	}

	assert.Equal(t, 3, list.Count())

	pots := list.GetPots()

	assert.Equal(t, int64(1000), pots[0].Level)
	assert.Equal(t, int64(2000), pots[1].Level)
	assert.Equal(t, int64(3000), pots[2].Level)

	assert.Equal(t, int64(1000), pots[0].Wager)
	assert.Equal(t, int64(1000), pots[1].Wager)
	assert.Equal(t, int64(1000), pots[2].Wager)

	prevPotWager := int64(0)
	for _, p := range pots {
		assert.Greater(t, p.Level, prevPotWager)
	}

	assert.Equal(t, 6, len(pots[0].Contributors))
	assert.Equal(t, 3, len(pots[1].Contributors))
	assert.Equal(t, 1, len(pots[2].Contributors))

	assert.Equal(t, int64(6000), pots[0].Total)
	assert.Equal(t, int64(3000), pots[1].Total)
	assert.Equal(t, int64(1000), pots[2].Total)
}

func TestPotList_SB_And_BB(t *testing.T) {

	contributers := []int64{
		5,
		10,
	}

	list := NewPotList()

	for idx, wager := range contributers {
		list.AddContributor(wager, idx)
	}

	assert.Equal(t, 2, list.Count())

	pots := list.GetPots()

	assert.Equal(t, int64(5), pots[0].Level)
	assert.Equal(t, int64(10), pots[1].Level)

	assert.Equal(t, 2, len(pots[0].Contributors))
	assert.Equal(t, 1, len(pots[1].Contributors))
}
