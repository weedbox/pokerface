package table

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/weedbox/pokerface"
	"github.com/weedbox/syncsaga"
)

var (
	ErrInvalidAction = errors.New("game: invalid action")
	ErrNoRunningGame = errors.New("game: no running game")
)

type Game interface {
	Start() error
	GetState() *pokerface.GameState
	OnStateUpdated(func(*pokerface.GameState))

	// Shortcut
	ReadyForAll() error
	PayAnte() error
	PayBlinds() error

	// Actions
	Ready(playerIdx int) error
	Pass(playerIdx int) error
	Pay(playerIdx int, chips int64) error
	Fold(playerIdx int) error
	Check(playerIdx int) error
	Call(playerIdx int) error
	Allin(playerIdx int) error
	Bet(playerIdx int, chips int64) error
	Raise(playerIdx int, chipLevel int64) error
}

type game struct {
	backend        Backend
	gs             *pokerface.GameState
	opts           *pokerface.GameOptions
	rg             *syncsaga.ReadyGroup
	mu             sync.RWMutex
	isClosed       bool
	incomingStates chan *pokerface.GameState
	onStateUpdated func(*pokerface.GameState)
}

func NewGame(backend Backend, opts *pokerface.GameOptions) *game {

	g := &game{
		backend:        backend,
		opts:           opts,
		rg:             syncsaga.NewReadyGroup(),
		incomingStates: make(chan *pokerface.GameState, 1024),
	}

	return g
}

func (g *game) runStateUpdater() {

	go func() {
		for state := range g.incomingStates {
			g.handleState(state)
		}
	}()
}

func (g *game) handleState(gs *pokerface.GameState) {

	switch gs.Status.CurrentEvent {
	case "GameClosed":
		g.Close()
	case "RoundClosed":

		// Next round automatically
		gs, err := g.backend.Next(gs)
		if err != nil {
			fmt.Println(err)
			return
		}

		g.updateState(gs)

	case "ReadyRequested":

		// Preparing ready group to wait for all player ready
		g.rg.Stop()
		g.rg.OnCompleted(func(rg *syncsaga.ReadyGroup) {
			g.ReadyForAll()
		})

		g.rg.ResetParticipants()
		for _, p := range gs.Players {
			g.rg.Add(int64(p.Idx), false)

			// Allow "ready" action
			p.AllowAction("ready")
		}

		g.rg.Start()

	case "AnteRequested":

		if gs.Meta.Ante == 0 {
			break
		}

		// Preparing ready group to wait for ante paid from all player
		g.rg.Stop()
		g.rg.OnCompleted(func(rg *syncsaga.ReadyGroup) {
			g.PayAnte()
		})

		g.rg.ResetParticipants()
		for _, p := range gs.Players {
			g.rg.Add(int64(p.Idx), false)

			// Allow "pay" action
			p.AllowAction("pay")
		}

		g.rg.Start()

	case "BlindsRequested":

		// Preparing ready group to wait for blinds
		g.rg.Stop()
		g.rg.OnCompleted(func(rg *syncsaga.ReadyGroup) {
			g.PayBlinds()
		})

		g.rg.ResetParticipants()
		for _, p := range gs.Players {
			if gs.Meta.Blind.BB > 0 && gs.HasPosition(p.Idx, "bb") {
				g.rg.Add(int64(p.Idx), false)
			} else if gs.Meta.Blind.SB > 0 && gs.HasPosition(p.Idx, "sb") {
				g.rg.Add(int64(p.Idx), false)
			} else if gs.Meta.Blind.Dealer > 0 && gs.HasPosition(p.Idx, "dealer") {
				g.rg.Add(int64(p.Idx), false)
			} else {
				continue
			}

			// Allow "pay" action
			p.AllowAction("pay")
		}

		g.rg.Start()
	}

	//fmt.Println("Game Updated =>", g.gs.Status.CurrentEvent)

	g.onStateUpdated(gs)
}

func (g *game) OnStateUpdated(fn func(*pokerface.GameState)) {
	g.onStateUpdated = fn
}

func (g *game) Start() error {

	g.runStateUpdater()

	gs, err := g.backend.CreateGame(g.opts)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) cloneState(gs *pokerface.GameState) *pokerface.GameState {

	// clone table state
	data, err := json.Marshal(gs)
	if err != nil {
		return nil
	}

	var state pokerface.GameState
	json.Unmarshal(data, &state)

	return &state
}

func (g *game) updateState(gs *pokerface.GameState) {

	g.mu.RLock()
	defer g.mu.RUnlock()

	state := g.cloneState(gs)
	g.gs = state

	if g.isClosed {
		return
	}

	g.incomingStates <- state

	//fmt.Println("Game Updating =>", g.gs.Status.CurrentEvent)
}

func (g *game) GetState() *pokerface.GameState {
	return g.gs
}

func (g *game) Close() {
	if g.isClosed {
		return
	}

	g.isClosed = true
	close(g.incomingStates)
}

func (g *game) Ready(playerIdx int) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "ready") || g.rg == nil {
		return ErrInvalidAction
	}

	//	fmt.Println("RRR", playerIdx)

	g.rg.Ready(int64(playerIdx))

	return nil
}

// Shortcut
func (g *game) ReadyForAll() error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	gs, err := g.backend.ReadyForAll(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) PayAnte() error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	gs, err := g.backend.PayAnte(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) PayBlinds() error {

	gs, err := g.backend.PayBlinds(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Pass(playerIdx int) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "pass") {
		return ErrInvalidAction
	}

	gs, err := g.backend.Pass(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Pay(playerIdx int, chips int64) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "pay") {
		return ErrInvalidAction
	}

	// For blinds
	switch g.gs.Status.CurrentEvent {
	case "AnteRequested":
		fallthrough
	case "BlindsRequested":
		g.rg.Ready(int64(playerIdx))
		return nil
	}

	gs, err := g.backend.Pay(g.gs, chips)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Fold(playerIdx int) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "fold") {
		return ErrInvalidAction
	}

	gs, err := g.backend.Fold(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Check(playerIdx int) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "check") {
		return ErrInvalidAction
	}

	gs, err := g.backend.Check(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Call(playerIdx int) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "call") {
		return ErrInvalidAction
	}

	gs, err := g.backend.Call(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Allin(playerIdx int) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "allin") {
		return ErrInvalidAction
	}

	gs, err := g.backend.Allin(g.gs)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Bet(playerIdx int, chips int64) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "bet") {
		return ErrInvalidAction
	}

	gs, err := g.backend.Bet(g.gs, chips)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}

func (g *game) Raise(playerIdx int, chipLevel int64) error {

	if g.gs == nil {
		return ErrNoRunningGame
	}

	p := g.gs.GetPlayer(playerIdx)
	if p == nil {
		return ErrPlayerNotInGame
	}

	if !g.gs.HasAction(playerIdx, "raise") {
		return ErrInvalidAction
	}

	gs, err := g.backend.Raise(g.gs, chipLevel)
	if err != nil {
		return err
	}

	g.updateState(gs)

	return nil
}
