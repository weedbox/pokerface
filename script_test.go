package pockerface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicCase(t *testing.T) {

	pf := NewPokerFace()

	opts := NewStardardGameOptions()

	// Preparing deck
	opts.Deck = NewStandardDeckCards()

	// Preparing players
	players := []*PlayerSetting{
		&PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"dealer"},
		},
		&PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"sb"},
		},
		&PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"bb"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)
	err := g.Start()
	assert.Nil(t, err)

	// Waiting for ready
	for _, p := range g.GetState().Players {
		assert.Equal(t, "Initialized", g.GetState().Status.CurrentEvent.Name)
		assert.Equal(t, 0, len(p.HoleCards))
		assert.Equal(t, false, p.Fold)
		assert.Equal(t, int64(0), p.Wager)
		assert.Equal(t, int64(0), p.Pot)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.Bankroll)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.InitialStackSize)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.StackSize)
		assert.Equal(t, "ready", p.AllowedActions[0])

		// Position checks
		if p.Idx == 0 {
			assert.Equal(t, "dealer", p.Positions[0])
		} else if p.Idx == 1 {
			assert.Equal(t, "sb", p.Positions[0])
		} else if p.Idx == 2 {
			assert.Equal(t, "bb", p.Positions[0])
		}

		err := g.Player(p.Idx).Ready()
		assert.Nil(t, err)
	}

	// Entering Preflop
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	g.PrintState()

	// Blinds
	for _, p := range g.GetState().Players {
		assert.Equal(t, 2, len(p.HoleCards))
		assert.Equal(t, 0, p.ActionCount)
		assert.Equal(t, false, p.Fold)
		assert.Equal(t, int64(0), p.Wager)
		assert.Equal(t, int64(0), p.Pot)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.Bankroll)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.InitialStackSize)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.StackSize)

		if p.Idx == 1 {
			// Small blind
			err := g.Player(p.Idx).Pay(5)
			assert.Nil(t, err)
		} else if p.Idx == 2 {
			// Big blind
			err := g.Player(p.Idx).Pay(10)
			assert.Nil(t, err)
		}
	}
}
