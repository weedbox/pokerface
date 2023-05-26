package pokerface

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cfsghost/pokerface/pot"
	"github.com/cfsghost/pokerface/task"
)

var (
	ErrNoDeck                      = errors.New("game: no deck")
	ErrNotEnoughBackroll           = errors.New("game: backroll is not enough")
	ErrInsufficientNumberOfPlayers = errors.New("game: insufficient number of players")
	ErrUnknownRound                = errors.New("game: unknown round")
	ErrNotFoundDealer              = errors.New("game: not found dealer")
	ErrUnknownTask                 = errors.New("game: Unknown task")
)

type Game interface {
	ApplyOptions(opts *GameOptions) error
	Start() error
	Resume() error
	GetEvent() *Event
	GetState() *GameState
	GetStateJSON() ([]byte, error)
	LoadState(gs *GameState) error
	Player(idx int) Player
	Deal(count int) []string
	Burn(count int) error
	FindDealer() *PlayerState
	ResetAllPlayerStatus() error
	StartAtDealer() (*PlayerState, error)
	NextMovablePlayer() *PlayerState
	SetCurrentPlayer(p *PlayerState) error
	GetCurrentPlayer() *PlayerState
	GetAllowedActions(player *PlayerState) []string
	GetAllowedBetActions(player *PlayerState) []string
	EmitEvent(event GameEvent, payload *EventPayload) error
	PrintState() error
}

type game struct {
	gs *GameState
}

func NewGame(opts *GameOptions) *game {
	g := &game{}
	g.ApplyOptions(opts)
	return g
}

func NewGameFromState(gs *GameState) *game {
	g := &game{}
	g.LoadState(gs)
	return g
}

func (g *game) GetState() *GameState {
	return g.gs
}

func (g *game) GetStateJSON() ([]byte, error) {
	return json.Marshal(g.gs)
}

func (g *game) LoadState(gs *GameState) error {
	g.gs = gs
	return g.Resume()
}

func (g *game) Resume() error {

	// emit event if state has event
	if g.gs.Status.CurrentEvent != nil {
		event := GameEventBySymbol[g.gs.Status.CurrentEvent.Name]

		fmt.Printf("Resume: %s\n", g.gs.Status.CurrentEvent.Name)

		// Activate by the last event
		g.EmitEvent(event, g.gs.Status.CurrentEvent.Payload)
	}

	return nil
}

func (g *game) ApplyOptions(opts *GameOptions) error {

	g.gs = &GameState{
		Players: make([]*PlayerState, 0),
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

	// Loading players
	for idx, p := range opts.Players {
		g.AddPlayer(idx, p)
	}

	return nil
}

func (g *game) AddPlayer(idx int, setting *PlayerSetting) error {

	ps := &PlayerState{
		Idx:              idx,
		Positions:        setting.Positions,
		Bankroll:         setting.Bankroll,
		InitialStackSize: setting.Bankroll,
		StackSize:        setting.Bankroll,
	}

	g.gs.Players = append(g.gs.Players, ps)

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
		state: p,
	}
}

func (g *game) Dealer() Player {
	for _, ps := range g.gs.Players {
		p := g.Player(ps.Idx)
		if p.CheckPosition("dealer") {
			return p
		}
	}

	return nil
}

func (g *game) SmallBlind() Player {
	for _, ps := range g.gs.Players {
		p := g.Player(ps.Idx)
		if p.CheckPosition("sb") {
			return p
		}
	}

	return nil
}

func (g *game) BigBlind() Player {
	for _, ps := range g.gs.Players {
		p := g.Player(ps.Idx)
		if p.CheckPosition("bb") {
			return p
		}
	}

	return nil
}

func (g *game) Deal(count int) []string {

	cards := make([]string, 0, count)

	finalPos := g.gs.Status.CurrentDeckPosition + count
	for i := g.gs.Status.CurrentDeckPosition; i < finalPos; i++ {
		cards = append(cards, g.gs.Meta.Deck[i])
		g.gs.Status.CurrentDeckPosition++
	}

	return cards
}

func (g *game) Burn(count int) error {
	g.gs.Status.Burned = append(g.gs.Status.Burned, g.Deal(count)...)
	return nil
}

func (g *game) FindDealer() *PlayerState {

	for _, p := range g.gs.Players {
		for _, pos := range p.Positions {
			if pos == "dealer" {
				return p
			}
		}
	}

	return nil
}

func (g *game) ResetAllPlayerAllowedActions() error {
	for _, p := range g.gs.Players {
		g.Player(p.Idx).Reset()
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

func (g *game) ResetRoundStatus() error {
	g.gs.Status.CurrentRoundPot = 0
	g.gs.Status.CurrentWager = 0
	g.gs.Status.CurrentRaiser = g.Dealer().State().Idx
	g.gs.Status.CurrentPlayer = g.gs.Status.CurrentRaiser
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

func (g *game) GetCurrentPlayer() *PlayerState {
	return g.Player(g.gs.Status.CurrentPlayer).State()
}

func (g *game) NextPlayer() *PlayerState {

	cur := g.gs.Status.CurrentPlayer

	for i := 1; i < len(g.gs.Players); i++ {

		// Find the next player
		cur++

		// The end of player list
		if cur == len(g.gs.Players) {
			cur = 0
		}

		p := g.gs.Players[cur]
		return p
	}

	return nil
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
			return p
		}
	}

	return nil
}

func (g *game) GetPlayers() []*PlayerState {

	players := make([]*PlayerState, 0)

	// Getting player list that dealer is the first element of it
	cur := g.Dealer().State().Idx

	for i := 1; i < len(g.gs.Players); i++ {

		// Find the next player
		cur++

		// The end of player list
		if cur == len(g.gs.Players) {
			cur = 0
		}

		p := g.gs.Players[cur]

		players = append(players, p)
	}

	return players
}

func (g *game) setCurrentPlayer(p *PlayerState) error {

	// Clear status
	if p == nil {
		g.gs.Status.CurrentPlayer = -1
		return nil
	}

	// Deside player who can move
	g.gs.Status.CurrentPlayer = p.Idx

	return nil
}

func (g *game) SetCurrentPlayer(p *PlayerState) error {

	if g.gs.Status.CurrentPlayer != -1 {
		// Clear allowed actions of current player
		g.Player(g.gs.Status.CurrentPlayer).ResetAllowedActions()
	}

	err := g.setCurrentPlayer(p)
	if err != nil {
		return err
	}

	if p != nil {
		// Figure out actions that player can be allowed to take
		g.Player(p.Idx).AllowActions(g.GetAllowedActions(p))
	}

	return nil
}

func (g *game) AlivePlayerCount() int {

	aliveCount := len(g.gs.Players)

	for _, p := range g.gs.Players {
		if p.Fold {
			aliveCount--
		}
	}

	return aliveCount
}

func (g *game) PlayerLoop() error {

	aliveCount := g.AlivePlayerCount()

	// only one player left
	if aliveCount == 1 {
		return g.EmitEvent(GameEvent_RoundClosed, nil)
	}

	// check for last check of first player on this round
	cp := g.GetCurrentPlayer()
	if g.gs.Status.CurrentRaiser == cp.Idx && cp.DidAction == "check" {
		return g.EmitEvent(GameEvent_RoundClosed, nil)
	}

	// next player
	p := g.NextMovablePlayer()

	//fmt.Printf("===================== cur=%d, actionCount=%d, raiser=%d\n", p.Idx, p.ActionCount, g.gs.Status.CurrentRaiser)

	if p.ActionCount == 0 {
		g.SetCurrentPlayer(p)
		// finally, it should trigger PlayerDidAction event
	} else if p.Idx != g.gs.Status.CurrentRaiser {
		g.SetCurrentPlayer(p)
		// finally, it should trigger PlayerDidAction event
	} else {
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

func (g *game) Start() error {

	// Initializing game status
	g.gs.Status.Pots = make([]*pot.Pot, 0)
	g.gs.Status.CurrentEvent = &Event{}

	return g.EmitEvent(GameEvent_Started, nil)
}

func (g *game) Initialize() error {

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

	g.ResetRoundStatus()

	return g.EmitEvent(GameEvent_Initialized, nil)
}

func (g *game) Prepare() error {

	if !g.WaitForAllPlayersReady("ready") {
		return nil
	}

	return g.EmitEvent(GameEvent_Prepared, nil)
}

func (g *game) RequestAnte() error {
	//TODO: preparing task for ante request
	return g.EmitEvent(GameEvent_AnteRequested, nil)
}

func (g *game) EnterPreflopRound() error {
	return g.EmitEvent(GameEvent_PreflopRoundEntered, nil)
}

func (g *game) EnterFlopRound() error {
	return g.EmitEvent(GameEvent_FlopRoundEntered, nil)
}

func (g *game) InitializeRound() error {

	// Initializing for stages (Preflop, Flop, Turn and River)
	switch g.gs.Status.Round {
	case "preflop":

		// Deal cards to players
		for _, p := range g.gs.Players {
			p.HoleCards = g.Deal(g.gs.Meta.HoleCardsCount)
		}
	case "flop":

		g.Burn(1)

		// Deal 3 board cards
		g.gs.Status.Board = append(g.gs.Status.Board, g.Deal(3)...)

		// Start at dealer
		_, err := g.StartAtDealer()
		if err != nil {
			return err
		}

	case "turn":
		fallthrough
	case "riber":

		g.Burn(1)

		// Deal board card
		g.gs.Status.Board = append(g.gs.Status.Board, g.Deal(1)...)

		// Start at dealer
		_, err := g.StartAtDealer()
		if err != nil {
			return err
		}
	}

	// Calculate power of the best combination for each player
	err := g.UpdateCombinationOfAllPlayers()
	if err != nil {
		return err
	}

	return g.EmitEvent(GameEvent_RoundInitialized, nil)
}

func (g *game) PrepareRound() error {

	fmt.Printf("Preparing round: %s\n", g.gs.Status.Round)

	if g.gs.Status.Round == "preflop" {
		return g.PreparePreflopRound()
	}

	// Other stage: start at dealer
	if !g.WaitForAllPlayersReady("ready") {
		return nil
	}

	_, err := g.StartAtDealer()
	if err != nil {
		return err
	}

	return g.EmitEvent(GameEvent_RoundPrepared, nil)
}

func (g *game) PreparePreflopRound() error {

	if g.gs.Meta.Blind.Dealer > 0 || g.gs.Meta.Blind.SB > 0 || g.gs.Meta.Blind.BB > 0 {

		// Initializing for current event
		event := g.GetEvent()

		// Task 1: request dealer blind
		if g.gs.Meta.Blind.Dealer > 0 && event.Payload.Task.GetTask("db") == nil {
			t := task.NewWaitPay("db")
			event.Payload.Task.AddTask(t)
		}

		// Task 2: request small blind
		if g.gs.Meta.Blind.SB > 0 && event.Payload.Task.GetTask("sb") == nil {
			t := task.NewWaitPay("sb")
			event.Payload.Task.AddTask(t)
		}

		// Task 3: request big blind
		if g.gs.Meta.Blind.BB > 0 && event.Payload.Task.GetTask("bb") == nil {
			t := task.NewWaitPay("bb")
			event.Payload.Task.AddTask(t)
		}

		// Task 4: Waiting for ready
		g.AssertReadyTask("ready")

		// Execute and check task status
		event.Payload.Task.Execute()

		if !event.Payload.Task.IsCompleted() {

			fmt.Printf("Check blinds: dealer=%d, sb=%d, bb=%d\n", g.gs.Meta.Blind.Dealer, g.gs.Meta.Blind.SB, g.gs.Meta.Blind.BB)

			// Getting available task for the next action
			t := event.Payload.Task.GetAvailableTask()

			// task for dealer blind and getting ready
			switch t.GetName() {
			case "db":

				// Start at dealer
				p := g.Dealer()
				p.AllowActions([]string{
					"pay",
				})
				g.setCurrentPlayer(p.State())
			case "sb":
				p := g.SmallBlind()
				p.AllowActions([]string{
					"pay",
				})
				g.setCurrentPlayer(p.State())
			case "bb":
				p := g.BigBlind()
				p.AllowActions([]string{
					"pay",
				})
				g.setCurrentPlayer(p.State())
			default:
				return ErrUnknownTask
			}

			fmt.Println("Waiting for blinds...")

			return nil
		}

		g.ResetAllPlayerAllowedActions()

		// Find the last player who has paid
		var lp *PlayerState
		players := g.GetPlayers()
		for _, p := range players {
			if p.Wager > 0 {
				lp = p
			}
		}

		g.SetCurrentPlayer(lp)

		return g.EmitEvent(GameEvent_RoundPrepared, nil)
	}

	return nil
}

func (g *game) onRoundClosed() error {

	// Update pots
	err := g.updatePots()
	if err != nil {
		return err
	}

	g.ResetRoundStatus()
	g.ResetAllPlayerStatus()

	aliveCount := g.AlivePlayerCount()
	if aliveCount == 1 {
		// Game is completed
		return g.EmitEvent(GameEvent_GameCompleted, nil)
	}

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

	//Note: this task is not required because we done need player ranking
	//ranks := g.CalculatePlayersRanking()

	// Calculate results with ranks
	err := g.CalculateGameResults()
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

func (g *game) PrintState() error {

	data, err := g.GetStateJSON()
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
