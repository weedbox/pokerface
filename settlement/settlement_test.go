package settlement

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface/pot"
)

func TestSinglePot(t *testing.T) {

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
	}

	for _, total := range pots {
		r.AddPot(total, []*pot.Level{
			&pot.Level{
				Level:        2000,
				Wager:        2000,
				Total:        6000,
				Contributors: []int{0, 1, 2},
			},
		})
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.UpdateScore(0, 1000)
	r.UpdateScore(1, 900)
	r.UpdateScore(2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 1, len(r.Pots[0].Winners))

	// finally, chips of player
	assert.Equal(t, int64(14000), r.Players[0].Final)
	assert.Equal(t, int64(8000), r.Players[1].Final)
	assert.Equal(t, int64(8000), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(4000), r.Players[0].Changed)
	assert.Equal(t, int64(-2000), r.Players[1].Changed)
	assert.Equal(t, int64(-2000), r.Players[2].Changed)
}

func TestMultiplePots(t *testing.T) {

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

	contributors := []struct {
		wager int64
		fold  bool
	}{
		{
			wager: 3000,
			fold:  false,
		},
		{
			wager: 4000,
			fold:  false,
		},
		{
			wager: 5000,
			fold:  false,
		},
	}

	ll := pot.NewLevelList()
	for idx, c := range contributors {
		ll.AddContributor(c.wager, idx, c.fold)
	}

	pots := ll.GetPots()
	for _, p := range pots {
		r.AddPot(p.Total, p.Levels)
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.UpdateScore(0, 1000) // Winner
	r.UpdateScore(1, 900)
	r.UpdateScore(2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 1, len(r.Pots[0].Winners))
	assert.Equal(t, 1, len(r.Pots[1].Winners))

	// finally, chips of player
	assert.Equal(t, int64(16000), r.Players[0].Final)
	assert.Equal(t, int64(8000), r.Players[1].Final)
	assert.Equal(t, int64(6000), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(6000), r.Players[0].Changed)
	assert.Equal(t, int64(-2000), r.Players[1].Changed)
	assert.Equal(t, int64(-4000), r.Players[2].Changed)
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

	contributors := []struct {
		wager int64
		fold  bool
	}{
		{
			wager: 3000,
			fold:  false,
		},
		{
			wager: 4000,
			fold:  false,
		},
		{
			wager: 5000,
			fold:  false,
		},
	}

	ll := pot.NewLevelList()
	for idx, c := range contributors {
		ll.AddContributor(c.wager, idx, c.fold)
	}

	pots := ll.GetPots()
	for _, p := range pots {
		r.AddPot(p.Total, p.Levels)
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.UpdateScore(0, 1000)
	r.UpdateScore(1, 1000)
	r.UpdateScore(2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 2, len(r.Pots[0].Winners))
	assert.Equal(t, 1, len(r.Pots[1].Winners))

	// finally, chips of player
	assert.Equal(t, int64(11500), r.Players[0].Final)
	assert.Equal(t, int64(12500), r.Players[1].Final)
	assert.Equal(t, int64(6000), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(1500), r.Players[0].Changed)
	assert.Equal(t, int64(2500), r.Players[1].Changed)
	assert.Equal(t, int64(-4000), r.Players[2].Changed)
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

	contributors := []struct {
		wager int64
		fold  bool
	}{
		{
			wager: 1000,
			fold:  false,
		},
		{
			wager: 10000,
			fold:  false,
		},
		{
			wager: 10000,
			fold:  false,
		},
	}

	ll := pot.NewLevelList()
	for idx, c := range contributors {
		ll.AddContributor(c.wager, idx, c.fold)
	}

	pots := ll.GetPots()
	for _, p := range pots {
		r.AddPot(p.Total, p.Levels)
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.UpdateScore(0, 1000)
	r.UpdateScore(1, 900)
	r.UpdateScore(2, 800)

	r.Calculate()

	// Pot winner
	assert.Equal(t, 1, len(r.Pots[0].Winners))
	assert.Equal(t, 1, len(r.Pots[1].Winners))

	// finally, chips of player
	assert.Equal(t, int64(3000), r.Players[0].Final)
	assert.Equal(t, int64(18000), r.Players[1].Final)
	assert.Equal(t, int64(0), r.Players[2].Final)

	// changes
	assert.Equal(t, int64(2000), r.Players[0].Changed)
	assert.Equal(t, int64(8000), r.Players[1].Changed)
	assert.Equal(t, int64(-10000), r.Players[2].Changed)
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
		r.AddPot(total, []*pot.Level{
			&pot.Level{
				Level:        1111,
				Wager:        1111,
				Total:        3333,
				Contributors: []int{0, 1, 2},
			},
		})
	}

	assert.Equal(t, len(r.Pots), len(pots))

	// Add contributers (pot index, player index, power score)
	r.UpdateScore(0, 1000)
	r.UpdateScore(1, 1000)
	r.UpdateScore(2, 800)

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
