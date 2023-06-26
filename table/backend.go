package table

import "github.com/weedbox/pokerface"

type Backend interface {
	CreateGame(opts *pokerface.GameOptions) (*pokerface.GameState, error)
	Next(gs *pokerface.GameState) (*pokerface.GameState, error)
	ReadyForAll(gs *pokerface.GameState) (*pokerface.GameState, error)
	PayAnte(gs *pokerface.GameState) (*pokerface.GameState, error)
	PayBlinds(gs *pokerface.GameState) (*pokerface.GameState, error)

	// Actions
	Call(gs *pokerface.GameState) (*pokerface.GameState, error)
	Pass(gs *pokerface.GameState) (*pokerface.GameState, error)
	Fold(gs *pokerface.GameState) (*pokerface.GameState, error)
	Check(gs *pokerface.GameState) (*pokerface.GameState, error)
	Allin(gs *pokerface.GameState) (*pokerface.GameState, error)
	Bet(gs *pokerface.GameState, chips int64) (*pokerface.GameState, error)
	Raise(gs *pokerface.GameState, chipLevel int64) (*pokerface.GameState, error)
	Pay(gs *pokerface.GameState, chips int64) (*pokerface.GameState, error)
}
