package table

import (
	"encoding/json"

	"github.com/weedbox/pokerface"
)

type NativeBackend struct {
	engine pokerface.PokerFace
}

func NewNativeBackend() *NativeBackend {
	return &NativeBackend{
		engine: pokerface.NewPokerFace(),
	}
}

func cloneState(gs *pokerface.GameState) *pokerface.GameState {

	//Note: we must clone a new structure for preventing original data of game engine is modified outside.
	data, err := json.Marshal(gs)
	if err != nil {
		return nil
	}

	var state pokerface.GameState
	err = json.Unmarshal([]byte(data), &state)
	if err != nil {
		return nil
	}

	return &state
}

func (nb *NativeBackend) getState(g pokerface.Game) *pokerface.GameState {
	return cloneState(g.GetState())
}

func (nb *NativeBackend) CreateGame(opts *pokerface.GameOptions) (*pokerface.GameState, error) {

	// Initializing game
	g := nb.engine.NewGame(opts)
	err := g.Start()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Next(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))
	err := g.Next()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) ReadyForAll(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))
	err := g.ReadyForAll()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Pass(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Pass()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) PayAnte(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.PayAnte()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) PayBlinds(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.PayBlinds()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Pay(gs *pokerface.GameState, chips int64) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Pay(chips)
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Fold(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Fold()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Check(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Check()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Call(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Call()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Allin(gs *pokerface.GameState) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Allin()
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Bet(gs *pokerface.GameState, chips int64) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Bet(chips)
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}

func (nb *NativeBackend) Raise(gs *pokerface.GameState, chipLevel int64) (*pokerface.GameState, error) {

	g := nb.engine.NewGameFromState(cloneState(gs))

	err := g.Raise(chipLevel)
	if err != nil {
		return nil, err
	}

	return nb.getState(g), nil
}
