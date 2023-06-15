package pokerface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface"
)

func Test_LastAction(t *testing.T) {

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

	// Start the game
	assert.Nil(t, g.Start())

	// Waiting for initial ready
	assert.Nil(t, g.ReadyForAll())

	// Ante
	assert.Equal(t, "Prepared", g.GetState().Status.CurrentEvent.Name)
	assert.Nil(t, g.PayAnte())
	assert.Equal(t, 8, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "ante", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	// Blinds
	assert.Nil(t, g.Pay(5))
	assert.Equal(t, 1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "small_blind", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(5), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Pay(10))
	assert.Equal(t, 2, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "big_blind", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	// Round: Preflop
	assert.Nil(t, g.ReadyForAll()) // ready for the round

	assert.Nil(t, g.Call())
	assert.Equal(t, 3, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 4, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 5, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 6, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 7, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 8, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // Dealer
	assert.Equal(t, 0, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(10), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // SB
	assert.Equal(t, 1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(5), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check()) // BB
	assert.Equal(t, 2, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	// Round: Flop
	assert.Nil(t, g.Next())
	assert.Equal(t, -1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "next", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.ReadyForAll()) // ready for the round

	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Check()) // SB
	assert.Equal(t, 1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check()) // BB
	assert.Equal(t, 2, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	// Bet from here
	assert.Nil(t, g.Bet(100))
	assert.Equal(t, 3, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "bet", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 4, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 5, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 6, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 7, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 8, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // Dealer
	assert.Equal(t, 0, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // SB
	assert.Equal(t, 1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // BB
	assert.Equal(t, 2, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	// Round: Turn
	assert.Nil(t, g.Next())
	assert.Equal(t, -1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "next", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.ReadyForAll()) // ready for the round

	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Check()) // SB
	assert.Equal(t, 1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Bet(100)) // BB
	assert.Equal(t, 2, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "bet", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Raise(200))
	assert.Equal(t, 3, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "raise", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(200), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Raise(300))
	assert.Equal(t, 4, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "raise", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(300), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 5, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(300), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 6, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(300), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 7, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(300), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 8, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(300), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // Dealer
	assert.Equal(t, 0, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(300), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // SB
	assert.Equal(t, 1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(300), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call()) // BB
	assert.Equal(t, 2, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(200), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Call())
	assert.Equal(t, 3, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "call", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(100), g.GetState().Status.LastAction.Value)

	// Round: River
	assert.Nil(t, g.Next())
	assert.Equal(t, -1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "next", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.ReadyForAll()) // ready for the round

	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Check()) // SB
	assert.Equal(t, 1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check()) // BB
	assert.Equal(t, 2, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check())
	assert.Equal(t, 3, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check())
	assert.Equal(t, 4, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check())
	assert.Equal(t, 5, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check())
	assert.Equal(t, 6, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check())
	assert.Equal(t, 7, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check())
	assert.Equal(t, 8, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	assert.Nil(t, g.Check()) // Dealer
	assert.Equal(t, 0, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "check", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)

	// Game closed
	assert.Nil(t, g.Next())
	assert.Equal(t, -1, g.GetState().Status.LastAction.Source)
	assert.Equal(t, "next", g.GetState().Status.LastAction.Type)
	assert.Equal(t, int64(0), g.GetState().Status.LastAction.Value)
}
