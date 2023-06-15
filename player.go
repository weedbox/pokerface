package pokerface

import (
	"errors"

	"github.com/weedbox/pokerface/task"
)

var (
	ErrInvalidAction = errors.New("player: invalid action")
	ErrIllegalRaise  = errors.New("player: illegal raise")
)

type Player interface {
	State() *PlayerState
	SeatIndex() int
	CheckPosition(pos string) bool
	AllowActions(actions []string) error
	ResetAllowedActions() error
	Reset() error
	Ready() error
	Pass() error
	Pay(chips int64) error
	Fold() error
	Check() error
	Call() error
	Allin() error
	Bet(chips int64) error
	Raise(chipLevel int64) error
}

type player struct {
	idx   int
	game  Game
	state *PlayerState
}

func (p *player) State() *PlayerState {

	state := p.game.GetState()

	if len(state.Players) <= p.idx {
		return nil
	}

	return state.Players[p.idx]
}

func (p *player) SeatIndex() int {
	return p.idx
}

func (p *player) Reset() error {
	p.state.Acted = false
	return p.ResetAllowedActions()
}

func (p *player) AllowActions(actions []string) error {
	p.state.AllowedActions = actions
	return nil
}

func (p *player) ResetAllowedActions() error {
	p.state.AllowedActions = make([]string, 0)
	return nil
}

func (p *player) IsMovable() bool {

	if p.game.GetState().Status.CurrentPlayer == p.idx {
		return true
	}

	if len(p.state.AllowedActions) > 0 {
		return true
	}

	return false
}

func (p *player) CheckPosition(pos string) bool {

	for _, p := range p.state.Positions {
		if p == pos {
			return true
		}
	}

	return false
}

func (p *player) CheckAction(action string) bool {

	for _, aa := range p.state.AllowedActions {
		if aa == action {
			return true
		}
	}

	return false
}

func (p *player) Ready() error {

	if !p.CheckAction("ready") {
		return ErrInvalidAction
	}

	//fmt.Printf("[Player %d] Get ready\n", p.idx)

	event := p.game.GetEvent()

	// Getting current task
	t := event.Payload.Task.GetAvailableTask()
	if t.GetType() != "ready" {
		return nil
	}

	// Update state to be ready
	wr := t.(*task.WaitReady)
	wr.Ready(p.idx)

	// Keep going
	return p.game.Resume()
}

func (p *player) Pass() error {

	if !p.CheckAction("pass") {
		return nil
	}

	p.state.Acted = true

	p.game.UpdateLastAction(p.idx, "pass", 0)

	return p.game.Resume()
}

func (p *player) pay(chips int64) error {

	if p.state.StackSize <= chips {

		// Update pot of current round
		gs := p.game.GetState()
		gs.Status.CurrentRoundPot += p.state.InitialStackSize - p.state.Wager

		p.state.DidAction = "allin"
		p.state.Wager = p.state.InitialStackSize
		p.state.StackSize = 0

		if p.state.InitialStackSize > gs.Status.CurrentWager {
			gs.Status.CurrentWager = p.state.InitialStackSize

			// Become new raiser
			p.game.BecomeRaiser(p)
		}

		return nil
	}

	p.state.Wager += chips
	p.state.StackSize = p.state.InitialStackSize - p.state.Wager

	// Update pot of current round
	gs := p.game.GetState()
	gs.Status.CurrentRoundPot += chips

	// player raised
	if gs.Status.CurrentWager < p.state.Wager {
		gs.Status.CurrentWager = p.state.Wager

		// Become new raiser
		p.game.BecomeRaiser(p)
	}

	return nil
}

func (p *player) Pay(chips int64) error {

	if !p.CheckAction("pay") {
		return ErrInvalidAction
	}

	//fmt.Printf("[Player %d] Pay %d\n", p.idx, chips)

	event := p.game.GetEvent()

	// Getting current task
	t := event.Payload.Task.GetAvailableTask()
	if t.GetType() != "pay" {
		return nil
	}

	// pay for wager
	err := p.pay(chips)
	if err != nil {
		return err
	}

	// Update task state
	wp := t.(*task.WaitPay)
	wp.Pay(p.idx, chips)

	gs := p.game.GetState()
	if gs.Status.CurrentEvent.Name == "Prepared" {
		p.game.UpdateLastAction(p.idx, "ante", chips)
	} else if gs.Status.CurrentEvent.Name == "RoundInitialized" {

		// Pay for blinds
		if p.CheckPosition("dealer") {
			p.game.UpdateLastAction(p.idx, "dealer_blind", chips)
		} else if p.CheckPosition("sb") {
			p.game.UpdateLastAction(p.idx, "small_blind", chips)
		} else {
			p.game.UpdateLastAction(p.idx, "big_blind", chips)
		}
	} else {
		p.game.UpdateLastAction(p.idx, "pay", chips)
	}

	// Keep going
	return p.game.Resume()
}

func (p *player) Fold() error {

	if !p.CheckAction("fold") {
		return ErrInvalidAction
	}

	p.state.Fold = true

	p.state.DidAction = "fold"
	p.state.Acted = true

	p.game.UpdateLastAction(p.idx, "fold", 0)

	return p.game.Resume()
}

func (p *player) Call() error {

	if !p.CheckAction("call") {
		return ErrInvalidAction
	}

	//fmt.Printf("[Player %d] call\n", p.idx)

	gs := p.game.GetState()

	delta := gs.Status.CurrentWager - p.state.Wager

	p.state.DidAction = "call"
	p.state.Acted = true

	p.pay(delta)

	p.game.UpdateLastAction(p.idx, "call", delta)

	return p.game.Resume()
}

func (p *player) Check() error {

	if !p.CheckAction("check") {
		return ErrInvalidAction
	}

	//fmt.Printf("[Player %d] check\n", p.idx)

	p.state.DidAction = "check"
	p.state.Acted = true

	p.game.UpdateLastAction(p.idx, "check", 0)

	return p.game.Resume()
}

func (p *player) Bet(chips int64) error {

	if !p.CheckAction("bet") {
		return ErrInvalidAction
	}

	//fmt.Printf("[Player %d] bet %d\n", p.idx, chips)

	p.state.DidAction = "bet"
	p.state.Acted = true

	p.pay(chips)

	p.game.GetState().Status.PreviousRaiseSize = chips

	p.game.UpdateLastAction(p.idx, "bet", chips)

	return p.game.Resume()
}

func (p *player) Raise(chipLevel int64) error {

	if !p.CheckAction("raise") {
		return ErrInvalidAction
	}

	gs := p.game.GetState()
	if chipLevel == 0 || chipLevel < gs.Status.CurrentWager {
		return ErrIllegalRaise
	}

	if chipLevel == gs.Status.CurrentWager {
		return p.Call()
	}

	// if chips is not enough to raise, player can do allin only
	raised := chipLevel - gs.Status.CurrentWager
	required := chipLevel - p.state.Wager
	//fmt.Println(gs.Status.PreviousRaiseSize)
	//fmt.Printf(" %d => initial=%d, raised=%d, required=%d\n", chipLevel, p.state.InitialStackSize, raised, required)
	if chipLevel >= p.state.InitialStackSize || raised < gs.Status.PreviousRaiseSize {
		return p.Allin()
	}

	//fmt.Printf("[Player %d] raise\n", p.idx)

	p.state.DidAction = "raise"
	p.state.Acted = true

	// Update raise size
	gs.Status.PreviousRaiseSize = chipLevel - gs.Status.CurrentWager

	p.pay(required)

	p.game.UpdateLastAction(p.idx, "raise", required)

	return p.game.Resume()
}

func (p *player) Allin() error {

	if !p.CheckAction("allin") {
		return ErrInvalidAction
	}

	//fmt.Printf("[Player %d] allin\n", p.idx)

	p.state.DidAction = "allin"
	p.state.Acted = true

	p.pay(p.state.StackSize)

	p.game.UpdateLastAction(p.idx, "allin", p.state.InitialStackSize)

	return p.game.Resume()
}
