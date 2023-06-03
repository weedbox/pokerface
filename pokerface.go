package pokerface

import (
	"time"

	"github.com/google/uuid"
)

type PokerFace interface {
	NewGame(opts *GameOptions) Game
	NewGameFromState(gs *GameState) Game
}

type pokerface struct {
}

func NewPokerFace() PokerFace {
	return &pokerface{}
}

func (pf *pokerface) NewGame(opts *GameOptions) Game {
	g := NewGame(opts)
	s := g.GetState()
	s.GameID = uuid.New().String()
	s.CreatedAt = time.Now().Unix()

	return g
}

func (pf *pokerface) NewGameFromState(gs *GameState) Game {
	return NewGameFromState(gs)
}
