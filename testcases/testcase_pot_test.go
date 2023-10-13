package pokerface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface"
)

func Test_Pot_Basic(t *testing.T) {

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
