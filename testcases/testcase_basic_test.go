package pokerface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface"
)

func PrepareAnte(t *testing.T, g pokerface.Game) {
	for _, p := range g.GetPlayers() {
		err := p.Pay(g.GetState().Meta.Ante)
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
	assert.Nil(t, g.Start())

	// Waiting for ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)

		// Position checks
		if p.SeatIndex() == 0 {
			assert.True(t, p.CheckPosition("dealer"))
		} else if p.SeatIndex() == 1 {
			assert.True(t, p.CheckPosition("sb"))
		} else if p.SeatIndex() == 2 {
			assert.True(t, p.CheckPosition("bb"))
		}
	}

	// Waiting for ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	//g.PrintState()

	// ante
	assert.Equal(t, "AnteRequested", g.GetState().Status.CurrentEvent)

	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)
	}

	assert.Nil(t, g.PayAnte())

	// Entering Preflop
	t.Log("Entering \"Prflop\" round")
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	assert.Equal(t, "BlindsRequested", g.GetState().Status.CurrentEvent)
	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, 2, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(10), p.State().Pot)
	}

	assert.Nil(t, g.PayBlinds())

	// Current wager on the table should be equal to big blind
	assert.Equal(t, int64(10), g.GetState().Status.CurrentWager)
	assert.Equal(t, 2, g.GetState().Status.CurrentRaiser)

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Starting player loop
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)

	// Dealer
	cp := g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// Dealer: call
	assert.Nil(t, cp.Call())

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// SB: call
	assert.Nil(t, cp.Call())

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 3, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "check", cp.State().AllowedActions[1])
	assert.Equal(t, "raise", cp.State().AllowedActions[2])

	// SB: check
	assert.Nil(t, cp.Check())

	// Entering Flop
	assert.Nil(t, g.Next())

	t.Log("Entering \"Flop\" round")
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Equal(t, "flop", g.GetState().Status.Round)
	assert.Equal(t, int64(30+3*g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Total)
	assert.Equal(t, int64(10+g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Wager)
	assert.Equal(t, 3, len(g.GetState().Status.Pots[0].Contributors))

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Starting player loop
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)

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
	assert.Nil(t, cp.Check())

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
	assert.Nil(t, cp.Check())

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
	assert.Nil(t, cp.Check())

	// Entering Turn
	assert.Nil(t, g.Next())

	t.Log("Entering \"Turn\" round")
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Equal(t, "turn", g.GetState().Status.Round)

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Starting player loop
	t.Log("Round is ready")
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)

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
	assert.Nil(t, cp.Check())

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
	assert.Nil(t, cp.Bet(30))

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
	assert.Nil(t, cp.Raise(60))

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
	assert.Nil(t, cp.Call())

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
	assert.Nil(t, cp.Call())

	// Entering River
	assert.Nil(t, g.Next())

	t.Log("Entering \"River\" round")
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Equal(t, "river", g.GetState().Status.Round)
	assert.Equal(t, int64(210+3*g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Total)
	assert.Equal(t, int64(70+g.GetState().Meta.Ante), g.GetState().Status.Pots[0].Wager)
	assert.Equal(t, 3, len(g.GetState().Status.Pots[0].Contributors))

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// SB
	cp = g.GetCurrentPlayer()
	assert.Nil(t, cp.Check())

	// BB
	cp = g.GetCurrentPlayer()
	assert.Nil(t, cp.Check())

	// Dealer
	cp = g.GetCurrentPlayer()
	assert.Nil(t, cp.Check())

	// Next
	assert.Nil(t, g.Next())
	assert.Equal(t, "GameClosed", g.GetState().Status.CurrentEvent)

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
	assert.Nil(t, g.Start())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Ante
	//PrepareAnte(t, g)
	assert.Nil(t, g.PayAnte())

	// Blinds
	assert.Nil(t, g.PayBlinds())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

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
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll())
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
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll())
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
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll())
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
	assert.Nil(t, g.Next())
}

func Test_Basic_SidePotCheck(t *testing.T) {

	pf := pokerface.NewPokerFace()

	// Options
	opts := pokerface.NewStardardGameOptions()
	opts.Blind.SB = 100
	opts.Blind.BB = 200

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
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)
	assert.Nil(t, g.Start())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Ante
	//assert.Nil(t, g.PayAnte())

	// Blinds
	assert.Nil(t, g.PayBlinds())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Preflop
	g.GetCurrentPlayer().Raise(800)
	g.GetCurrentPlayer().Call() // Dealer
	g.GetCurrentPlayer().Fold() // SB
	g.GetCurrentPlayer().Fold() // BB

	// No side pot out there
	assert.Equal(t, 1, len(g.GetState().Status.Pots))

	// Flop
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll())
	g.GetCurrentPlayer().Pass()  // SB
	g.GetCurrentPlayer().Pass()  // BB
	g.GetCurrentPlayer().Check() // UG
	g.GetCurrentPlayer().Check() // Dealer

	// Turn
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll())
	g.GetCurrentPlayer().Pass()  // SB
	g.GetCurrentPlayer().Pass()  // BB
	g.GetCurrentPlayer().Check() // UG
	g.GetCurrentPlayer().Check() // Dealer

	// River
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll())
	g.GetCurrentPlayer().Pass()  // SB
	g.GetCurrentPlayer().Pass()  // BB
	g.GetCurrentPlayer().Check() // UG
	g.GetCurrentPlayer().Check() // Dealer

	// Game closed
	assert.Nil(t, g.Next())
}
