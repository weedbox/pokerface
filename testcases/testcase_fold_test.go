package pokerface

import (
	"testing"

	"github.com/cfsghost/pokerface"
	"github.com/stretchr/testify/assert"
)

func Test_Fold(t *testing.T) {

	pf := pokerface.NewPokerFace()

	opts := pokerface.NewStardardGameOptions()
	opts.Ante = 10

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"dealer"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"sb"},
		},
		&pokerface.PlayerSetting{
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
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "Initialized", g.GetState().Status.CurrentEvent.Name)
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)
		assert.Equal(t, "ready", p.State().AllowedActions[0])

		// Position checks
		if p.SeatIndex() == 0 {
			assert.True(t, p.CheckPosition("dealer"))
		} else if p.SeatIndex() == 1 {
			assert.True(t, p.CheckPosition("sb"))
		} else if p.SeatIndex() == 2 {
			assert.True(t, p.CheckPosition("bb"))
		}

		err := p.Ready()
		assert.Nil(t, err)
	}

	// ante
	assert.Equal(t, "Prepared", g.GetState().Status.CurrentEvent.Name)

	for _, p := range g.GetPlayers() {
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, 0, p.State().ActionCount)
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)
		assert.Equal(t, "pay", p.State().AllowedActions[0])
		err := p.Pay(opts.Ante)
		assert.Nil(t, err)
	}

	// Entering Preflop
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
		assert.Equal(t, 2, len(p.State().HoleCards))
		assert.Equal(t, 0, p.State().ActionCount)
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(10), p.State().Pot)

		if p.SeatIndex() == 1 {
			assert.Equal(t, "pay", p.State().AllowedActions[0])

			// Small blind
			err := p.Pay(5)
			assert.Nil(t, err)
		} else if p.SeatIndex() == 2 {
			assert.Equal(t, "pay", p.State().AllowedActions[0])

			// Big blind
			err := p.Pay(10)
			assert.Nil(t, err)
		}
	}

	// Current wager on the table should be equal to big blind
	assert.Equal(t, int64(10), g.GetState().Status.CurrentWager)
	assert.Equal(t, 2, g.GetState().Status.CurrentRaiser)

	// get ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "ready", p.State().AllowedActions[0])
		err := p.Ready()
		assert.Nil(t, err)
	}

	// Starting player loop
	assert.Equal(t, "RoundPrepared", g.GetState().Status.CurrentEvent.Name)

	// Dealer
	cp := g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// Dealer: fold
	err = g.Player(cp.SeatIndex()).Fold()
	assert.Nil(t, err)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// SB: fold
	err = g.Player(cp.SeatIndex()).Fold()
	assert.Nil(t, err)

	// This game should be closed immediately
	assert.Equal(t, "GameClosed", g.GetState().Status.CurrentEvent.Name)

	//g.PrintState()
}
