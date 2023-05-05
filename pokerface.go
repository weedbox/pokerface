package main

type PokerFace interface {
	NewGame() Game
}

type pokerface struct {
}

func NewPokerFace() PokerFace {
	return &pokerface{}
}

func (pf *pokerface) NewGame() Game {
	return NewGame()
}
