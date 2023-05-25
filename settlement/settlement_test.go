package settlement

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {

	r := NewResult()

	// Bankroll of players
	players := []int64{
		10000,
		10000,
		10000,
	}

	for idx, bankroll := range players {
		r.AddPlayer(idx, bankroll)
	}

	// Pot
	pots := []int64{
		6000,
		3000,
	}

	for _, total := range pots {
		r.AddPot(total)
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.AddContributor(0, 0, 1000)
	r.AddContributor(0, 1, 900)
	r.AddContributor(0, 2, 800)
	r.AddContributor(1, 0, 1000)
	r.AddContributor(1, 1, 900)
	r.AddContributor(1, 2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 1, len(r.Pots[0].Winners))
	assert.Equal(t, 1, len(r.Pots[1].Winners))

	// finally, chips of player
	assert.Equal(t, int64(16000), r.Players[0].Final)
	assert.Equal(t, int64(7000), r.Players[1].Final)
	assert.Equal(t, int64(7000), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(6000), r.Players[0].Changed)
	assert.Equal(t, int64(-3000), r.Players[1].Changed)
	assert.Equal(t, int64(-3000), r.Players[2].Changed)
}

func TestMultipleWinners(t *testing.T) {

	r := NewResult()

	// Bankroll of players
	players := []int64{
		10000,
		10000,
		10000,
	}

	for idx, bankroll := range players {
		r.AddPlayer(idx, bankroll)
	}

	// Pot
	pots := []int64{
		6000,
		3000,
	}

	for _, total := range pots {
		r.AddPot(total)
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.AddContributor(0, 0, 1000)
	r.AddContributor(0, 1, 1000)
	r.AddContributor(0, 2, 800)
	r.AddContributor(1, 0, 1000)
	r.AddContributor(1, 1, 1000)
	r.AddContributor(1, 2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 2, len(r.Pots[0].Winners))
	assert.Equal(t, 2, len(r.Pots[1].Winners))

	// finally, chips of player
	assert.Equal(t, int64(11500), r.Players[0].Final)
	assert.Equal(t, int64(11500), r.Players[1].Final)
	assert.Equal(t, int64(7000), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(1500), r.Players[0].Changed)
	assert.Equal(t, int64(1500), r.Players[1].Changed)
	assert.Equal(t, int64(-3000), r.Players[2].Changed)
}

func TestAllin(t *testing.T) {

	r := NewResult()

	// Bankroll of players
	players := []int64{
		1000,
		10000,
		10000,
	}

	for idx, bankroll := range players {
		r.AddPlayer(idx, bankroll)
	}

	// Pot
	pots := []int64{
		3000,
		2000,
	}

	for _, total := range pots {
		r.AddPot(total)
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.AddContributor(0, 0, 1000)
	r.AddContributor(0, 1, 900)
	r.AddContributor(0, 2, 800)
	r.AddContributor(1, 1, 900)
	r.AddContributor(1, 2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 1, len(r.Pots[0].Winners))
	assert.Equal(t, 1, len(r.Pots[1].Winners))

	// finally, chips of player
	assert.Equal(t, int64(3000), r.Players[0].Final)
	assert.Equal(t, int64(10000), r.Players[1].Final)
	assert.Equal(t, int64(8000), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(2000), r.Players[0].Changed)
	assert.Equal(t, int64(0), r.Players[1].Changed)
	assert.Equal(t, int64(-2000), r.Players[2].Changed)
}

func TestMultipleWinnersWithRemainder(t *testing.T) {

	r := NewResult()

	// Bankroll of players
	players := []int64{
		10000,
		10000,
		10000,
	}

	for idx, bankroll := range players {
		r.AddPlayer(idx, bankroll)
	}

	// Pot
	pots := []int64{
		3333,
	}

	for _, total := range pots {
		r.AddPot(total)
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.AddContributor(0, 0, 1000)
	r.AddContributor(0, 1, 1000)
	r.AddContributor(0, 2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 2, len(r.Pots[0].Winners))

	// finally, chips of player
	assert.Equal(t, int64(10556), r.Players[0].Final)
	assert.Equal(t, int64(10555), r.Players[1].Final)
	assert.Equal(t, int64(8889), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(556), r.Players[0].Changed)
	assert.Equal(t, int64(555), r.Players[1].Changed)
	assert.Equal(t, int64(-1111), r.Players[2].Changed)
}
