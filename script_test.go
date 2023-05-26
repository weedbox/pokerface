package pokerface

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

	//TODO: ante

	// Entering Preflop
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	for _, p := range g.GetState().Players {
		assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
		assert.Equal(t, 2, len(p.HoleCards))
		assert.Equal(t, 0, p.ActionCount)
		assert.Equal(t, false, p.Fold)
		assert.Equal(t, int64(0), p.Wager)
		assert.Equal(t, int64(0), p.Pot)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.Bankroll)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.InitialStackSize)
		assert.Equal(t, int64(players[p.Idx].Bankroll), p.StackSize)

		if p.Idx == 1 {
			assert.Equal(t, "pay", p.AllowedActions[0])

			// Small blind
			err := g.Player(p.Idx).Pay(5)
			assert.Nil(t, err)
		} else if p.Idx == 2 {
			assert.Equal(t, "pay", p.AllowedActions[0])

			// Big blind
			err := g.Player(p.Idx).Pay(10)
			assert.Nil(t, err)
		}
	}

	// Current wager on the table should be equal to big blind
	assert.Equal(t, int64(10), g.GetState().Status.CurrentWager)
	assert.Equal(t, 2, g.GetState().Status.CurrentRaiser)

	// get ready
	for _, p := range g.GetState().Players {
		assert.Equal(t, "ready", p.AllowedActions[0])
		err := g.Player(p.Idx).Ready()
		assert.Nil(t, err)
	}

	// Starting player loop
	assert.Equal(t, "RoundPrepared", g.GetState().Status.CurrentEvent.Name)

	// Dealer
	cp := g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "fold", cp.AllowedActions[1])
	assert.Equal(t, "call", cp.AllowedActions[2])
	assert.Equal(t, "raise", cp.AllowedActions[3])

	// Dealer: call
	err = g.Player(cp.Idx).Call()
	assert.Nil(t, err)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "fold", cp.AllowedActions[1])
	assert.Equal(t, "call", cp.AllowedActions[2])
	assert.Equal(t, "raise", cp.AllowedActions[3])

	// SB: call
	err = g.Player(cp.Idx).Call()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 3, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "check", cp.AllowedActions[1])
	assert.Equal(t, "raise", cp.AllowedActions[2])

	// SB: check
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	// Entering Flop
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "flop", g.GetState().Status.Round)
	assert.Equal(t, int64(30), g.GetState().Status.Pots[0].Total)
	assert.Equal(t, int64(10), g.GetState().Status.Pots[0].Wager)
	assert.Equal(t, 3, len(g.GetState().Status.Pots[0].Contributors))

	// get ready
	for _, p := range g.GetState().Players {
		assert.Equal(t, "ready", p.AllowedActions[0])
		err := g.Player(p.Idx).Ready()
		assert.Nil(t, err)
	}

	// Starting player loop
	assert.Equal(t, "RoundPrepared", g.GetState().Status.CurrentEvent.Name)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 3, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "check", cp.AllowedActions[1])
	assert.Equal(t, "bet", cp.AllowedActions[2])

	// SB: check
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 3, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "check", cp.AllowedActions[1])
	assert.Equal(t, "bet", cp.AllowedActions[2])

	// BB: check
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	// Dealer
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 3, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "check", cp.AllowedActions[1])
	assert.Equal(t, "bet", cp.AllowedActions[2])

	// Dealer: check
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	// Entering Turn
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "turn", g.GetState().Status.Round)

	// get ready
	for _, p := range g.GetState().Players {
		assert.Equal(t, "ready", p.AllowedActions[0])
		err := g.Player(p.Idx).Ready()
		assert.Nil(t, err)
	}

	// Starting player loop
	assert.Equal(t, "RoundPrepared", g.GetState().Status.CurrentEvent.Name)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 3, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "check", cp.AllowedActions[1])
	assert.Equal(t, "bet", cp.AllowedActions[2])

	// SB: check
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 3, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "check", cp.AllowedActions[1])
	assert.Equal(t, "bet", cp.AllowedActions[2])

	// BB: bet 30
	err = g.Player(cp.Idx).Bet(30)
	assert.Nil(t, err)

	// Dealer
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 4, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "fold", cp.AllowedActions[1])
	assert.Equal(t, "call", cp.AllowedActions[2])
	assert.Equal(t, "raise", cp.AllowedActions[3])

	// Dealer: raise 60
	err = g.Player(cp.Idx).Raise(60)
	assert.Nil(t, err)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 4, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "fold", cp.AllowedActions[1])
	assert.Equal(t, "call", cp.AllowedActions[2])
	assert.Equal(t, "raise", cp.AllowedActions[3])

	// SB: call
	err = g.Player(cp.Idx).Call()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(30), cp.Wager)
	assert.Equal(t, cp.InitialStackSize, cp.Bankroll-cp.Pot)
	assert.Equal(t, cp.StackSize, cp.Bankroll-cp.Pot-cp.Wager)
	assert.Equal(t, 4, len(cp.AllowedActions))
	assert.Equal(t, "allin", cp.AllowedActions[0])
	assert.Equal(t, "fold", cp.AllowedActions[1])
	assert.Equal(t, "call", cp.AllowedActions[2])
	assert.Equal(t, "raise", cp.AllowedActions[3])

	// BB: call
	err = g.Player(cp.Idx).Call()
	assert.Nil(t, err)

	// Entering River
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "river", g.GetState().Status.Round)
	assert.Equal(t, int64(210), g.GetState().Status.Pots[0].Total)
	assert.Equal(t, int64(70), g.GetState().Status.Pots[0].Wager)
	assert.Equal(t, 3, len(g.GetState().Status.Pots[0].Contributors))

	// get ready
	for _, p := range g.GetState().Players {
		assert.Equal(t, "ready", p.AllowedActions[0])
		err := g.Player(p.Idx).Ready()
		assert.Nil(t, err)
	}

	// SB
	cp = g.GetCurrentPlayer()
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	// Dealer
	cp = g.GetCurrentPlayer()
	err = g.Player(cp.Idx).Check()
	assert.Nil(t, err)

	g.PrintState()
}
