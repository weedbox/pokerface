package main

import (
	"fmt"

	"github.com/cfsghost/pokerface/waitgroup"
)

type Player interface {
	State() *PlayerState

	CheckPosition(pos string) bool
	AllowActions(actions []string) error
	Ready() error
	Pass() error
	Pay(chips int64) error
	Fold() error
	Check() error
	Bet(chips int64) error
	Raise(chips int64) error
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

func (p *player) AllowActions(actions []string) error {
	p.state.AllowedActions = actions
	return nil
}

func (p *player) IsMovable() bool {

	if p.game.GetState().Status.CurrentPlayer == p.idx {
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
		return nil
	}

	fmt.Printf("[Player %d] Get ready\n", p.idx)

	wg := p.game.GetWaitGroup()
	if wg == nil {
		return nil
	}

	// Check waitgroup type
	if wg.Type != waitgroup.TypeReady {
		return nil
	}

	wg.GetStateByIdx(p.idx).State = true

	p.game.Resume()

	return nil
}

func (p *player) Pass() error {

	if !p.CheckAction("pass") {
		return nil
	}

	// Implement the logic for the Pass() function
	return nil
}

func (p *player) pay(chips int64) error {

	if p.state.StackSize < chips {
		p.state.Wager += p.state.StackSize
		p.state.StackSize = 0
		return nil
	}

	p.state.Wager += chips
	p.state.StackSize += p.state.Wager

	return nil
}

func (p *player) PayAnte(chips int64) error {

	if !p.CheckAction("pay_ante") {
		return nil
	}

	fmt.Printf("[Player %d] pay ante %d\n", p.idx, chips)

	wg := p.game.GetWaitGroup()
	if wg == nil {
		return nil
	}

	// Check waitgroup type
	if wg.Type != waitgroup.TypePayAnte {
		return nil
	}

	p.pay(chips)

	wg.GetStateByIdx(p.idx).State = true

	p.game.Resume()

	return nil
}

func (p *player) Pay(chips int64) error {

	if !p.CheckAction("pay") {
		return nil
	}

	// Implement the logic for the Pay() function

	return nil
}

func (p *player) Fold() error {

	if !p.CheckAction("fold") {
		return nil
	}

	// Implement the logic for the Fold() function
	return nil
}

func (p *player) Check() error {

	if !p.CheckAction("check") {
		return nil
	}

	// Implement the logic for the Check() function
	return nil
}

func (p *player) Bet(chips int64) error {

	if !p.CheckAction("bet") {
		return nil
	}

	// Implement the logic for the Bet() function
	return nil
}

func (p *player) Raise(chips int64) error {

	if !p.CheckAction("raise") {
		return nil
	}

	// Implement the logic for the Raise() function
	return nil
}
