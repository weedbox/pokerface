package pokerface

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/weedbox/pokerface/pot"
	"github.com/weedbox/pokerface/task"
)

var (
	ErrNoDeck                      = errors.New("game: no deck")
	ErrNotEnoughBackroll           = errors.New("game: backroll is not enough")
	ErrNoDealer                    = errors.New("game: no dealer")
	ErrInsufficientNumberOfPlayers = errors.New("game: insufficient number of players")
	ErrUnknownRound                = errors.New("game: unknown round")
	ErrNotFoundDealer              = errors.New("game: not found dealer")
	ErrUnknownTask                 = errors.New("game: unknown task")
	ErrNotClosedRound              = errors.New("game: round is not closed")
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
	Dealer() Player
	SmallBlind() Player
	BigBlind() Player
	Deal(count int) []string
	Burn(count int) error
	BecomeRaiser(Player) error
	ResetAllPlayerStatus() error
	StartAtDealer() (Player, error)
	GetPlayerCount() int
	GetPlayers() []Player
	SetCurrentPlayer(Player) error
	GetCurrentPlayer() Player
	GetAllowedActions(Player) []string
	GetAvailableActions(Player) []string
	GetAlivePlayerCount() int
	GetMovablePlayerCount() int
	Next() error
	EmitEvent(event GameEvent, payload *EventPayload) error
	PrintState() error

	// Actions
	ReadyForAll() error
	Pass() error
	Ready(playerIdx int) error
	PayAnte() error
	Pay(chips int64) error
	Fold() error
	Check() error
	Call() error
	Allin() error
	Bet(chips int64) error
	Raise(chipLevel int64) error
}

type game struct {
	gs         *GameState
	players    map[int]Player
	dealer     Player
	smallBlind Player
	bigBlind   Player
}

func NewGame(opts *GameOptions) *game {
	g := &game{
		players: make(map[int]Player),
	}
	g.ApplyOptions(opts)
	return g
}

func NewGameFromState(gs *GameState) *game {
	g := &game{
		players: make(map[int]Player),
	}
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

	// Initializing players
	for _, ps := range g.gs.Players {
		g.addPlayer(ps)
	}

	return g.Resume()
}

func (g *game) Resume() error {

	// emit event if state has event
	if g.gs.Status.CurrentEvent != nil {
		event := GameEventBySymbol[g.gs.Status.CurrentEvent.Name]

		//fmt.Printf("Resume: %s\n", g.gs.Status.CurrentEvent.Name)

		// Activate by the last event
		return g.EmitEvent(event, g.gs.Status.CurrentEvent.Payload)
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

func (g *game) addPlayer(state *PlayerState) error {

	// Create player instance
	p := &player{
		idx:   state.Idx,
		game:  g,
		state: state,
	}

	if p.CheckPosition("dealer") {
		g.dealer = p
	}

	if p.CheckPosition("sb") {
		g.smallBlind = p
	} else if p.CheckPosition("bb") {
		g.bigBlind = p
	}

	g.players[state.Idx] = p

	return nil
}

func (g *game) AddPlayer(idx int, setting *PlayerSetting) error {

	// Create player state
	ps := &PlayerState{
		Idx:              idx,
		Positions:        setting.Positions,
		Bankroll:         setting.Bankroll,
		InitialStackSize: setting.Bankroll,
		StackSize:        setting.Bankroll,
		Combination:      &CombinationInfo{},
	}

	g.gs.Players = append(g.gs.Players, ps)

	return g.addPlayer(ps)
}

func (g *game) Player(idx int) Player {

	if idx < 0 || idx >= g.GetPlayerCount() {
		return nil
	}

	return g.players[idx]
}

func (g *game) Dealer() Player {
	return g.dealer
}

func (g *game) SmallBlind() Player {
	return g.smallBlind
}

func (g *game) BigBlind() Player {
	return g.bigBlind
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

func (g *game) ResetAllPlayerAllowedActions() error {
	for _, p := range g.GetPlayers() {
		p.Reset()
	}

	return nil
}

func (g *game) ResetAllPlayerStatus() error {
	for _, p := range g.GetPlayers() {
		ps := p.State()
		ps.AllowedActions = make([]string, 0)
		ps.Pot += ps.Wager
		ps.Wager = 0
		ps.InitialStackSize = ps.StackSize

		if ps.Fold {
			ps.DidAction = "fold"
		} else if ps.InitialStackSize == 0 {
			ps.DidAction = "allin"
		} else {
			ps.DidAction = ""
		}
	}

	return nil
}

func (g *game) ResetRoundStatus() error {
	g.gs.Status.PreviousRaiseSize = 0
	g.gs.Status.CurrentRoundPot = 0
	g.gs.Status.CurrentWager = 0
	g.gs.Status.CurrentRaiser = g.Dealer().State().Idx
	g.gs.Status.CurrentPlayer = g.gs.Status.CurrentRaiser
	return nil
}

func (g *game) StartAtDealer() (Player, error) {

	// Start at dealer
	dealer := g.Dealer()
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

func (g *game) GetCurrentPlayer() Player {
	return g.Player(g.gs.Status.CurrentPlayer)
}

func (g *game) NextPlayer() Player {

	cur := g.gs.Status.CurrentPlayer
	playerCount := g.GetPlayerCount()

	for i := 1; i < playerCount; i++ {

		// Find the next player
		cur++

		// The end of player list
		if cur == playerCount {
			cur = 0
		}

		p := g.gs.Players[cur]
		return g.Player(p.Idx)
	}

	return nil
}

func (g *game) GetPlayerCount() int {
	return len(g.gs.Players)
}

func (g *game) GetPlayers() []Player {

	players := make([]Player, 0)
	playerCount := g.GetPlayerCount()

	// Getting player list that dealer should be the first element of it
	cur := g.Dealer().SeatIndex()

	for i := 0; i < playerCount; i++ {

		players = append(players, g.players[cur])

		// Find the next player
		cur++

		// The end of player list
		if cur == playerCount {
			cur = 0
		}
	}

	return players
}

func (g *game) setCurrentPlayer(p Player) error {

	// Clear status
	if p == nil {
		g.gs.Status.CurrentPlayer = -1
		return nil
	}

	// Deside player who can move
	g.gs.Status.CurrentPlayer = p.SeatIndex()

	return nil
}

func (g *game) SetCurrentPlayer(p Player) error {

	if g.gs.Status.CurrentPlayer != -1 {
		// Clear allowed actions of current player
		g.GetCurrentPlayer().ResetAllowedActions()
	}

	err := g.setCurrentPlayer(p)
	if err != nil {
		return err
	}

	if p != nil {
		// Figure out actions that player can be allowed to take
		actions := g.GetAllowedActions(p)
		p.AllowActions(actions)
	}

	return nil
}

func (g *game) GetAlivePlayerCount() int {

	aliveCount := g.GetPlayerCount()

	for _, p := range g.gs.Players {
		if p.Fold {
			aliveCount--
		}
	}

	return aliveCount
}

func (g *game) GetMovablePlayerCount() int {

	mCount := g.GetPlayerCount()

	for _, p := range g.gs.Players {
		// Fold or allin
		if p.Fold || p.StackSize == 0 {
			mCount--
		}
	}

	return mCount
}

func (g *game) BecomeRaiser(p Player) error {

	if p.State().Wager > 0 {
		p.State().VPIP = true
	}

	g.gs.Status.CurrentRaiser = p.SeatIndex()

	// Reset all player states except raiser
	for _, ps := range g.gs.Players {
		if ps.Idx == p.SeatIndex() {
			continue
		}

		if ps.DidAction != "fold" && ps.DidAction != "allin" {
			ps.DidAction = ""
		}

		ps.Acted = false
	}

	return nil
}

func (g *game) RequestPlayerAction() error {

	// only one player left
	if g.GetAlivePlayerCount() == 1 {
		return g.EmitEvent(GameEvent_RoundClosed, nil)
	}

	// no player can move because everybody did all-in already for this game
	if g.GetMovablePlayerCount() == 0 {
		return g.EmitEvent(GameEvent_RoundClosed, nil)
	}

	// next player
	p := g.NextPlayer()

	//fmt.Printf("===================== [%s] cur=%d, actionCount=%d, raiser=%d\n", g.gs.Status.Round, p.SeatIndex(), p.State().ActionCount, g.gs.Status.CurrentRaiser)

	// Run around already, no one need to act
	if p.State().Acted {
		return g.EmitEvent(GameEvent_RoundClosed, nil)
	}

	return g.SetCurrentPlayer(p)
}

func (g *game) GetAllowedActions(p Player) []string {

	// player is movable for this round
	if g.gs.Status.CurrentPlayer == p.SeatIndex() {
		return g.GetAvailableActions(p)
	}

	return make([]string, 0)
}

func (g *game) GetAvailableActions(p Player) []string {

	actions := make([]string, 0)

	// Invalid
	if p == nil {
		return actions
	}

	ps := p.State()

	if ps.Fold {
		actions = append(actions, "pass")
		return actions
	}

	// chips left
	if ps.StackSize == 0 {
		actions = append(actions, "pass")
		return actions
	} else {
		actions = append(actions, "allin")
	}

	if ps.Wager < g.gs.Status.CurrentWager {
		actions = append(actions, "fold")

		// call
		if ps.InitialStackSize > g.gs.Status.CurrentWager {

			actions = append(actions, "call")

			// raise
			if ps.InitialStackSize > g.gs.Status.CurrentWager+g.gs.Status.PreviousRaiseSize {
				actions = append(actions, "raise")
			}
		}

	} else {
		actions = append(actions, "check")

		if ps.InitialStackSize >= g.gs.Status.MiniBet {
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

	// Check the number of players
	if g.GetPlayerCount() < 2 {
		return ErrInsufficientNumberOfPlayers
	}

	// Require dealer
	if g.dealer == nil {
		return ErrNoDealer
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

	// Initializing game status
	g.gs.Status.Pots = make([]*pot.Pot, 0)
	g.gs.Status.Board = make([]string, 0)
	g.gs.Status.Burned = make([]string, 0)
	g.gs.Status.CurrentEvent = &Event{}

	return g.EmitEvent(GameEvent_Started, nil)
}

func (g *game) Initialize() error {

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

	if !g.WaitForPayment("ante") {
		return nil
	}

	g.ResetRoundStatus()

	return g.EmitEvent(GameEvent_AnteRequested, nil)

}

func (g *game) Next() error {

	switch g.gs.Status.Round {
	case "preflop":
		fallthrough
	case "flop":
		fallthrough
	case "turn":
		fallthrough
	case "river":
		return g.nextRound()
	}

	return nil
}

func (g *game) nextRound() error {

	if g.gs.Status.CurrentEvent.Name != "RoundClosed" {
		return ErrNotClosedRound
	}

	g.ResetRoundStatus()
	g.ResetAllPlayerStatus()

	if g.GetAlivePlayerCount() == 1 {
		// Game is completed
		return g.EmitEvent(GameEvent_GameCompleted, nil)
	}

	// Going to the next round
	switch g.gs.Status.Round {
	case "preflop":
		return g.EnterFlopRound()
	case "flop":
		return g.EnterTurnRound()
	case "turn":
		return g.EnterRiverRound()
	case "river":
		return g.EmitEvent(GameEvent_GameCompleted, nil)
	}

	return ErrUnknownRound
}

func (g *game) EnterPreflopRound() error {
	g.gs.Status.Round = "preflop"
	return g.EmitEvent(GameEvent_PreflopRoundEntered, nil)
}

func (g *game) EnterFlopRound() error {
	g.gs.Status.Round = "flop"
	return g.EmitEvent(GameEvent_FlopRoundEntered, nil)
}

func (g *game) EnterTurnRound() error {
	g.gs.Status.Round = "turn"
	return g.EmitEvent(GameEvent_TurnRoundEntered, nil)
}

func (g *game) EnterRiverRound() error {
	g.gs.Status.Round = "river"
	return g.EmitEvent(GameEvent_RiverRoundEntered, nil)
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
	case "river":

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

	//fmt.Printf("Preparing round: %s\n", g.gs.Status.Round)

	if g.gs.Status.Round == "preflop" {
		return g.PreparePreflopRound()
	}

	// Everybody did all-in, no need to keep going with normal way
	if g.GetMovablePlayerCount() == 0 {
		return g.EmitEvent(GameEvent_RoundClosed, nil)
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

	if g.gs.Meta.Blind.Dealer == 0 && g.gs.Meta.Blind.SB == 0 && g.gs.Meta.Blind.BB > 0 {
		return nil
	}

	// Initializing for current event
	event := g.GetEvent()

	// Task 1: request dealer blind
	if g.gs.Meta.Blind.Dealer > 0 && event.Payload.Task.GetTask("db") == nil {
		t := task.NewWaitPay("db", g.gs.Meta.Blind.Dealer)

		playerIdx := g.Dealer().State().Idx
		t.PrepareStates([]int{playerIdx})

		event.Payload.Task.AddTask(t)
	}

	// Task 2: request small blind
	if g.gs.Meta.Blind.SB > 0 && event.Payload.Task.GetTask("sb") == nil {
		t := task.NewWaitPay("sb", g.gs.Meta.Blind.SB)

		playerIdx := g.SmallBlind().State().Idx
		t.PrepareStates([]int{playerIdx})

		event.Payload.Task.AddTask(t)
	}

	// Task 3: request big blind
	if g.gs.Meta.Blind.BB > 0 && event.Payload.Task.GetTask("bb") == nil {

		// Minimal raise size
		g.gs.Status.PreviousRaiseSize = g.gs.Meta.Blind.BB

		t := task.NewWaitPay("bb", g.gs.Meta.Blind.BB)

		playerIdx := g.BigBlind().State().Idx
		t.PrepareStates([]int{playerIdx})

		event.Payload.Task.AddTask(t)
	}

	// Task 4: Waiting for ready
	g.AssertReadyTask("ready")

	// Execute and check task status
	event.Payload.Task.Execute()

	if !event.Payload.Task.IsCompleted() {

		//fmt.Printf("Check blinds: dealer=%d, sb=%d, bb=%d\n", g.gs.Meta.Blind.Dealer, g.gs.Meta.Blind.SB, g.gs.Meta.Blind.BB)

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
			g.setCurrentPlayer(p)
		case "sb":
			p := g.SmallBlind()
			p.AllowActions([]string{
				"pay",
			})
			g.setCurrentPlayer(p)
		case "bb":
			p := g.BigBlind()
			p.AllowActions([]string{
				"pay",
			})
			g.setCurrentPlayer(p)
		case "ready":
			// Do nothing
		default:
			return ErrUnknownTask
		}

		//fmt.Println("Waiting for blinds...")

		return nil
	}

	g.ResetAllPlayerAllowedActions()

	// Everybody did all-in, no need to keep going with normal way
	if g.GetMovablePlayerCount() == 0 {
		return g.EmitEvent(GameEvent_RoundClosed, nil)
	}

	// Find the last player who has paid
	var lp Player
	for i, p := range g.GetPlayers() {

		// First one is dealer, skip it
		if i == 0 {
			continue
		}

		if p.State().Wager == 0 {
			break
		}

		lp = p
	}

	g.SetCurrentPlayer(lp)

	return g.EmitEvent(GameEvent_RoundPrepared, nil)
}

func (g *game) PrintState() error {

	data, err := g.GetStateJSON()
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
