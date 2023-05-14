package main

import (
	"errors"
	"fmt"

	"github.com/cfsghost/pokerface/waitgroup"
)

type GameEvent int32

const (

	// Initialization
	GameEvent_Started GameEvent = iota
	GameEvent_Initialized
	GameEvent_AnteRequested
	GameEvent_AnteReceived

	// States
	GameEvent_Dealt
	GameEvent_WagerRequested

	// Rounds
	GameEvent_PreflopRoundEntered
	GameEvent_FlopRoundEntered
	GameEvent_TurnRoundEntered
	GameEvent_RiverRoundEntered
	GameEvent_RoundInitialized
	GameEvent_RoundPrepared
	GameEvent_RoundClosed
	GameEvent_PlayerDidAction

	// Result
	GameEvent_GameCompleted
	GameEvent_SettlementRequested
	GameEvent_SettlementCompleted
	GameEvent_GameClosed
)

var GameEventSymbols = []string{
	"Started",
	"Initialized",
	"AnteRequested",
	"AnteReceived",
	"Dealt",
	"WagerRequested",
	"PreflopRoundEntered",
	"FlopRoundEntered",
	"TurnRoundEntered",
	"RiverRoundEntered",
	"RoundInitialized",
	"RoundPrepared",
	"RoundClosed",
	"PlayerDidAction",
	"GameCompleted",
	"SettlementRequested",
	"SettlementCompleted",
	"GameClosed",
}

var GameEventBySymbol = map[string]GameEvent{
	"Started":             GameEvent_Started,
	"Initialized":         GameEvent_Initialized,
	"AnteRequested":       GameEvent_AnteRequested,
	"AnteReceived":        GameEvent_AnteReceived,
	"Dealt":               GameEvent_Dealt,
	"WagerRequested":      GameEvent_WagerRequested,
	"PreflopRoundEntered": GameEvent_PreflopRoundEntered,
	"FlopRoundEntered":    GameEvent_FlopRoundEntered,
	"TurnRoundEntered":    GameEvent_TurnRoundEntered,
	"RiverRoundEntered":   GameEvent_RiverRoundEntered,
	"RoundInitialized":    GameEvent_RoundInitialized,
	"RoundPrepared":       GameEvent_RoundPrepared,
	"RoundClosed":         GameEvent_RoundClosed,
	"PlayerDidAction":     GameEvent_PlayerDidAction,
	"GameCompleted":       GameEvent_GameCompleted,
	"SettlementRequested": GameEvent_SettlementRequested,
	"SettlementCompleted": GameEvent_SettlementCompleted,
	"GameClosed":          GameEvent_GameClosed,
}

var (
	ErrNoDeck                      = errors.New("game: no deck")
	ErrNotEnoughBackroll           = errors.New("game: backroll is not enough")
	ErrInsufficientNumberOfPlayers = errors.New("game: insufficient number of players")
	ErrUnknownRound                = errors.New("game: unknown round")
	ErrNotFoundDealer              = errors.New("game: not found dealer")
)

type Game interface {
	ApplyOptions(opts *GameOptions) error
	GetWaitGroup() *waitgroup.WaitGroup
	GetState() *GameState
	LoadState(gs *GameState) error
	Player(idx int) Player
	Deal(count int) []string
	Burn(count int) error
	FindDealer() *PlayerState
	ResetAllPlayerStatus() error
	StartAtDealer() (*PlayerState, error)
	NextMovablePlayer() *PlayerState
	SetCurrentPlayer(p *PlayerState) error
	GetAllowedActions(player *PlayerState) []string
	GetAllowedBetActions(player *PlayerState) []string
	EmitEvent(event GameEvent, runtime interface{}) error
}

type game struct {
	gs *GameState
	wg *waitgroup.WaitGroup
}

func NewGame() *game {
	return &game{}
}

func (g *game) GetWaitGroup() *waitgroup.WaitGroup {
	return g.gs
}

func (g *game) GetState() *GameState {
	return g.gs
}

func (g *game) LoadState(gs *GameState) error {
	g.gs = gs

	// emit event if state has event
	if g.gs.Status.CurrentEvent != nil {
		event := GameEventBySymbol[g.gs.Status.CurrentEvent.Name]

		// Activate by the last event
		g.EmitEvent(event, g.gs.Status.CurrentEvent.Runtime)
	}

	return nil
}

func (g *game) ApplyOptions(opts *GameOptions) error {

	g.gs = &GameState{
		Meta: Meta{
			Ante:                   opts.Ante,
			Blind:                  opts.Blind,
			Limit:                  opts.Limit,
			HoleCardsCount:         opts.HoleCardsCount,
			RequiredHoleCardsCount: opts.RequiredHoleCardsCount,
			CombinationPowers:      opts.CombinationPowers,
			Deck:                   opts.Deck,
			BurnCount:              opts.BurnCount,
		},
	}

	return nil
}

func (g *game) Player(idx int) Player {

	if idx < 0 || idx >= len(g.gs.Players) {
		return nil
	}

	p := g.gs.Players[idx]

	return &player{
		idx:   idx,
		game:  g,
		state: &p,
	}
}

func (g *game) Deal(count int) []string {

	cards := make([]string, 0, count)

	finalPos := g.gs.Status.CurrentDeckPosition + count
	for i := g.gs.Status.CurrentDeckPosition + 1; i <= finalPos; i++ {
		cards = append(cards, g.gs.Meta.Deck[i])
	}

	return cards
}

func (g *game) Burn(count int) error {
	g.gs.Status.Burned = append(g.gs.Status.Burned, g.Deal(count)...)
	return nil
}

func (g *game) FindDealer() *PlayerState {

	for _, p := range g.gs.Players {
		for _, pos := range p.Position {
			if pos == "dealer" {
				return &p
			}
		}
	}

	return nil
}

func (g *game) ResetAllPlayerStatus() error {
	for _, p := range g.gs.Players {
		p.DidAction = ""
		p.AllowedActions = make([]string, 0)
		p.Pot += p.Wager
		p.Wager = 0
		p.InitialStackSize = p.StackSize
	}

	return nil
}

func (g *game) StartAtDealer() (*PlayerState, error) {

	// Start at dealer
	dealer := g.FindDealer()
	if dealer == nil {
		return nil, ErrNotFoundDealer
	}

	// Update status
	err := g.SetCurrentPlayer(dealer)
	if err != nil {
		return nil, err
	}

	return dealer, nil
}

func (g *game) NextMovablePlayer() *PlayerState {

	cur := g.gs.Status.CurrentPlayer

	for i := 1; i < len(g.gs.Players); i++ {

		// Find the next player
		cur++

		// The end of player list
		if cur == len(g.gs.Players) {
			cur = 0
		}

		p := g.gs.Players[cur]

		// Find the player who did not fold
		if !p.Fold {
			return &p
		}
	}

	return nil
}

func (g *game) SetCurrentPlayer(p *PlayerState) error {

	// Clear status
	if p == nil {
		g.gs.Status.CurrentPlayer = -1
		return nil
	}

	// Deside player who can move
	g.gs.Status.CurrentPlayer = p.Idx

	// Figure out actions that player can be allowed to take
	g.Player(p.Idx).AllowActions(g.GetAllowedActions(p))

	return nil
}

func (g *game) PlayerLoop() error {

	liveCount := len(g.gs.Players)

	for _, p := range g.gs.Players {
		if p.Fold {
			liveCount--
		}
	}

	// only one player left
	if liveCount < 2 {
		// Game is completed
		return g.EmitEvent(GameEvent_GameCompleted, nil)
	}

	// next player
	p := g.NextMovablePlayer()

	if p.ActionCount == 0 {
		g.SetCurrentPlayer(p)
		// finally, it should trigger PlayerDidAction event
	} else if p.Idx != g.gs.Status.CurrentRaiser {
		g.SetCurrentPlayer(p)
		// finally, it should trigger PlayerDidAction event
	} else {
		g.SetCurrentPlayer(nil)
		// No more player can move
		return g.EmitEvent(GameEvent_RoundClosed, nil)
	}

	return nil
}

func (g *game) GetAllowedActions(player *PlayerState) []string {

	// player is movable for this round
	if g.gs.Status.CurrentPlayer == player.Idx {
		return g.GetAllowedBetActions(player)
	}

	return make([]string, 0)
}

func (g *game) GetAllowedBetActions(player *PlayerState) []string {

	actions := make([]string, 0)

	// Invalid player state
	if player == nil {
		return actions
	}

	if player.Fold {
		actions = append(actions, "pass")
		return actions
	}

	// chips left
	if player.StackSize == 0 {
		actions = append(actions, "pass")
		return actions
	} else {
		actions = append(actions, "allin")
	}

	if player.Wager < g.gs.Status.CurrentWager {
		actions = append(actions, "fold")

		// call
		if player.InitialStackSize > g.gs.Status.CurrentWager {

			actions = append(actions, "call")

			// raise
			if player.InitialStackSize >= g.gs.Status.CurrentWager*2 {
				actions = append(actions, "raise")
			}
		}

	} else {
		actions = append(actions, "check")

		if player.InitialStackSize >= g.gs.Status.MiniBet {
			if g.gs.Status.CurrentWager == 0 {
				actions = append(actions, "bet")
			} else {
				actions = append(actions, "raise")
			}
		}
	}

	return actions
}

func (g *game) triggerEvent(event GameEvent) error {

	switch event {

	case GameEvent_Started:
		fmt.Println("Game has started.")
		return g.onStarted()

	case GameEvent_Initialized:
		fmt.Println("Game has been initialized.")
		return g.onInitialized()
		//g.EmitEvent(GameEvent_Initialized)

	case GameEvent_AnteRequested:
		fmt.Println("Ante has been requested.")
		return g.onAnteRequested()

	case GameEvent_AnteReceived:
		fmt.Println("Ante is received.")
		return g.onAnteReceived()

	case GameEvent_Dealt:
		fmt.Println("Cards have been dealt.")
		//		return g.antePreparation()

	case GameEvent_WagerRequested:
		fmt.Println("Wager has been requested.")

	case GameEvent_PreflopRoundEntered:
		fmt.Println("Entered Preflop round.")
		return g.onPreflopRoundEntered()

	case GameEvent_FlopRoundEntered:
		fmt.Println("Entered Flop round.")
		return g.onFlopRoundEntered()

	case GameEvent_TurnRoundEntered:
		fmt.Println("Entered Turn round.")
		return g.onTurnRoundEntered()

	case GameEvent_RiverRoundEntered:
		fmt.Println("Entered River round.")
		return g.onRiverRoundEntered()

	case GameEvent_RoundInitialized:
		fmt.Println("Current round has initialized.")
		return g.onRoundInitialized()

	case GameEvent_RoundPrepared:
		fmt.Println("Current round has been prepared.")
		return g.onRoundPrepared()

	case GameEvent_RoundClosed:
		fmt.Println("Current round has closed.")
		return g.onRoundClosed()

	case GameEvent_PlayerDidAction:
		fmt.Println("Player did action.")
		return g.onPlayerDidAction()

	case GameEvent_GameCompleted:
		fmt.Println("Game has been completed.")
		return g.onGameCompleted()

	case GameEvent_SettlementRequested:
		fmt.Println("Settlement has been requested.")
		return g.onSettlementRequested()

	case GameEvent_SettlementCompleted:
		fmt.Println("Settlement has been completed.")
		return g.onSettlementCompleted()

	case GameEvent_GameClosed:
		fmt.Println("Game has closed.")
	}

	return nil
}

func (g *game) EmitEvent(event GameEvent, runtime interface{}) error {

	// Update current event
	g.gs.Status.CurrentEvent.Name = GameEventSymbols[event]
	g.gs.Status.CurrentEvent.Runtime = runtime

	return g.triggerEvent(event)
}

func (g *game) onStarted() error {

	// Check the number of players
	if len(g.gs.Players) < 2 {
		return ErrInsufficientNumberOfPlayers
	}

	// Check backroll
	for _, p := range g.gs.Players {

		if p.Bankroll <= 0 {
			return ErrNotEnoughBackroll
		}
	}

	// No desk was set
	if len(g.gs.Meta.Deck) == 0 {
		return ErrNoDeck
	}

	// Shuffle cards
	g.gs.Meta.Deck = ShuffleCards(g.gs.Meta.Deck)

	// Initialize minimum bet
	if g.gs.Meta.Blind.Dealer > g.gs.Meta.Blind.BB {
		g.gs.Status.MiniBet = g.gs.Meta.Blind.Dealer
	} else {
		g.gs.Status.MiniBet = g.gs.Meta.Blind.BB
	}

	return g.EmitEvent(GameEvent_Initialized, nil)
}

func (g *game) onInitialized() error {

	wg, err := g.WaitForAllPlayersReady()
	if err != nil {
		return err
	}

	if wg.IsCompleted() {
		g.wg = nil
		return g.EmitEvent(GameEvent_AnteRequested, nil)
	}

	// Nothing to do, just waiting for all players to be ready

	return nil
}

func (g *game) onAnteRequested() error {

	wg, err := g.WaitForAllPlayersPaidAnte()
	if err != nil {
		return err
	}

	if wg.IsCompleted() {
		g.wg = nil
		return g.EmitEvent(GameEvent_AnteReceived, nil)
	}

	// Nothing to do, just waiting for all players paid

	return nil
}

func (g *game) onAnteReceived() error {
	return g.EmitEvent(GameEvent_PreflopRoundEntered, nil)
}

func (g *game) onRoundInitialized() error {
	//return g.EmitEvent(GameEvent_RoundClosed)

	// Preparing round
	if g.gs.Meta.Blind.Dealer > 0 || g.gs.Meta.Blind.SB > 0 || g.gs.Meta.Blind.BB > 0 {
		//TODO: required players to pay for blinds
	} else {

		// Start at dealer
		_, err := g.StartAtDealer()
		if err != nil {
			return err
		}
	}

	return g.EmitEvent(GameEvent_RoundPrepared, nil)
}

func (g *game) onRoundPrepared() error {

	wg, err := g.WaitForAllPlayersReady()
	if err != nil {
		return err
	}

	if wg.IsCompleted() {
		g.wg = nil

		// All players is ready
		return g.PlayerLoop()
	}

	// Nothing to do, just waiting for all players to be ready

	return nil
}

func (g *game) onPreflopRoundEntered() error {

	g.gs.Status.Round = "preflop"

	// Deal cards to players
	for _, p := range g.gs.Players {
		p.HoleCards = g.Deal(g.gs.Meta.HoleCardsCount)
	}

	return g.EmitEvent(GameEvent_RoundInitialized, nil)
}

func (g *game) onFlopRoundEntered() error {

	g.gs.Status.Round = "flop"

	g.Burn(1)

	// Board
	g.gs.Status.Board = append(g.gs.Status.Board, g.Deal(3)...)

	// Start at dealer
	_, err := g.StartAtDealer()
	if err != nil {
		return err
	}

	return g.EmitEvent(GameEvent_RoundInitialized, nil)
}

func (g *game) onTurnRoundEntered() error {

	g.gs.Status.Round = "turn"

	g.Burn(1)

	// Board
	g.gs.Status.Board = append(g.gs.Status.Board, g.Deal(1)...)

	// Start at dealer
	_, err := g.StartAtDealer()
	if err != nil {
		return err
	}

	return g.EmitEvent(GameEvent_RoundInitialized, nil)
}

func (g *game) onRiverRoundEntered() error {

	g.gs.Status.Round = "river"

	g.Burn(1)

	// Board
	g.gs.Status.Board = append(g.gs.Status.Board, g.Deal(1)...)

	// Start at dealer
	_, err := g.StartAtDealer()
	if err != nil {
		return err
	}

	return g.EmitEvent(GameEvent_RoundInitialized, nil)
}

func (g *game) onRoundClosed() error {

	g.ResetAllPlayerStatus()

	switch g.gs.Status.Round {
	case "preflop":
		return g.EmitEvent(GameEvent_FlopRoundEntered, nil)
	case "flop":
		return g.EmitEvent(GameEvent_TurnRoundEntered, nil)
	case "turn":
		return g.EmitEvent(GameEvent_RiverRoundEntered, nil)
	case "river":
		return g.EmitEvent(GameEvent_GameCompleted, nil)
	}

	return ErrUnknownRound
}

func (g *game) onPlayerDidAction() error {
	return g.PlayerLoop()
}

func (g *game) onGameCompleted() error {
	return g.EmitEvent(GameEvent_SettlementRequested, nil)
}

func (g *game) onSettlementRequested() error {

	ranks := g.CalculatePowerRanking()

	// Calculate results with ranks
	err := g.CalculateGameResults(ranks)
	if err != nil {
		return err
	}

	return g.EmitEvent(GameEvent_SettlementCompleted, nil)
}

func (g *game) onSettlementCompleted() error {
	return g.EmitEvent(GameEvent_GameClosed, nil)
}

func (g *game) onGameClosed() error {
	return nil
}
