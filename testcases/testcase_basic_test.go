package pokerface

import (
	"testing"

	"github.com/cfsghost/pokerface"
	"github.com/stretchr/testify/assert"
)

func PrepareAnte(t *testing.T, g pokerface.Game) {
	for _, p := range g.GetPlayers() {
		err := p.Pay(g.GetState().Meta.Ante)
		assert.Nil(t, err)
	}
}

func AllPlayersReady(t *testing.T, g pokerface.Game) {
	for _, p := range g.GetPlayers() {
		err := p.Ready()
		assert.Nil(t, err)
	}
}

func Test_Basic(t *testing.T) {

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

	//g.PrintState()

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
	t.Log("Entering \"Prflop\" round")
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
	assert.Equal(t, "PlayerActionRequested", g.GetState().Status.CurrentEvent.Name)

	// Dealer
	cp := g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// Dealer: call
	err = cp.Call()
	assert.Nil(t, err)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// SB: call
	err = cp.Call()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 3, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "check", cp.State().AllowedActions[1])
	assert.Equal(t, "raise", cp.State().AllowedActions[2])

	// SB: check
	err = cp.Check()
	assert.Nil(t, err)

	// Entering Flop
	err = g.Next()
	assert.Nil(t, err)

	t.Log("Entering \"Flop\" round")
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "flop", g.GetState().Status.Round)
	assert.Equal(t, int64(30+3*g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Total)
	assert.Equal(t, int64(10+g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Wager)
	assert.Equal(t, 3, len(g.GetState().Status.Pots[0].Contributors))

	// get ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "ready", p.State().AllowedActions[0])
		err := p.Ready()
		assert.Nil(t, err)
	}

	// Starting player loop
	assert.Equal(t, "PlayerActionRequested", g.GetState().Status.CurrentEvent.Name)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 3, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "check", cp.State().AllowedActions[1])
	assert.Equal(t, "bet", cp.State().AllowedActions[2])

	// SB: check
	err = cp.Check()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 3, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "check", cp.State().AllowedActions[1])
	assert.Equal(t, "bet", cp.State().AllowedActions[2])

	// BB: check
	err = cp.Check()
	assert.Nil(t, err)

	// Dealer
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 3, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "check", cp.State().AllowedActions[1])
	assert.Equal(t, "bet", cp.State().AllowedActions[2])

	// Dealer: check
	err = cp.Check()
	assert.Nil(t, err)

	// Entering Turn
	err = g.Next()
	assert.Nil(t, err)

	t.Log("Entering \"Turn\" round")
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "turn", g.GetState().Status.Round)

	// get ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "ready", p.State().AllowedActions[0])
		err := p.Ready()
		assert.Nil(t, err)
	}

	// Starting player loop
	t.Log("Round is ready")
	assert.Equal(t, "PlayerActionRequested", g.GetState().Status.CurrentEvent.Name)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 3, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "check", cp.State().AllowedActions[1])
	assert.Equal(t, "bet", cp.State().AllowedActions[2])

	// SB: check
	err = cp.Check()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(0), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 3, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "check", cp.State().AllowedActions[1])
	assert.Equal(t, "bet", cp.State().AllowedActions[2])

	// BB: bet 30
	err = cp.Bet(30)
	assert.Nil(t, err)

	// Dealer
	cp = g.GetCurrentPlayer()
	assert.Equal(t, "dealer", cp.State().Positions[0])
	assert.Equal(t, int64(0), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// Dealer: raise to 60 (+30)
	err = cp.Raise(60)
	assert.Nil(t, err)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, "sb", cp.State().Positions[0])
	assert.Equal(t, int64(0), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// SB: call
	err = cp.Call()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, "bb", cp.State().Positions[0])
	assert.Equal(t, int64(30), cp.State().Wager)
	assert.Equal(t, cp.State().InitialStackSize, cp.State().Bankroll-cp.State().Pot)
	assert.Equal(t, cp.State().StackSize, cp.State().Bankroll-cp.State().Pot-cp.State().Wager)
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// BB: call
	err = cp.Call()
	assert.Nil(t, err)

	// Entering River
	err = g.Next()
	assert.Nil(t, err)

	t.Log("Entering \"River\" round")
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "river", g.GetState().Status.Round)
	assert.Equal(t, int64(210+3*g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Total)
	assert.Equal(t, int64(70+g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Wager)
	assert.Equal(t, 3, len(g.GetState().Status.Pots[0].Contributors))

	// get ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "ready", p.State().AllowedActions[0])
		err := p.Ready()
		assert.Nil(t, err)
	}

	// SB
	cp = g.GetCurrentPlayer()
	err = cp.Check()
	assert.Nil(t, err)

	// BB
	cp = g.GetCurrentPlayer()
	err = cp.Check()
	assert.Nil(t, err)

	// Dealer
	cp = g.GetCurrentPlayer()
	err = cp.Check()
	assert.Nil(t, err)

	// Next
	err = g.Next()
	assert.Nil(t, err)
	assert.Equal(t, "GameClosed", g.GetState().Status.CurrentEvent.Name)

	//g.PrintState()
}

func Test_Basic_NinePlayers(t *testing.T) {

	pf := pokerface.NewPokerFace()

	// Options
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
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)
	err := g.Start()
	assert.Nil(t, err)

	// Waiting for ready
	AllPlayersReady(t, g)

	// Ante
	PrepareAnte(t, g)

	// Blinds
	g.GetCurrentPlayer().Pay(5)
	g.GetCurrentPlayer().Pay(10)

	// Get ready for preflop
	AllPlayersReady(t, g)

	// Preflop
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()  // Dealer
	g.GetCurrentPlayer().Call()  // SB
	g.GetCurrentPlayer().Check() // BB

	// Flop
	err = g.Next()
	assert.Nil(t, err)
	AllPlayersReady(t, g)
	g.GetCurrentPlayer().Check() // SB
	g.GetCurrentPlayer().Check() // BB
	g.GetCurrentPlayer().Bet(100)
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call() // Dealer
	g.GetCurrentPlayer().Call() // SB
	g.GetCurrentPlayer().Call() // BB

	// Turn
	err = g.Next()
	assert.Nil(t, err)
	AllPlayersReady(t, g)
	g.GetCurrentPlayer().Check()  // SB
	g.GetCurrentPlayer().Bet(100) // BB
	g.GetCurrentPlayer().Raise(200)
	g.GetCurrentPlayer().Raise(300)
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call()
	g.GetCurrentPlayer().Call() // Dealer
	g.GetCurrentPlayer().Call() // SB
	g.GetCurrentPlayer().Call() // BB
	g.GetCurrentPlayer().Call()

	// River
	err = g.Next()
	assert.Nil(t, err)
	AllPlayersReady(t, g)
	g.GetCurrentPlayer().Check() // SB
	g.GetCurrentPlayer().Check() // BB
	g.GetCurrentPlayer().Check()
	g.GetCurrentPlayer().Check()
	g.GetCurrentPlayer().Check()
	g.GetCurrentPlayer().Check()
	g.GetCurrentPlayer().Check()
	g.GetCurrentPlayer().Check()
	g.GetCurrentPlayer().Check() // Dealer

	// Game closed
	err = g.Next()
	assert.Nil(t, err)
}
