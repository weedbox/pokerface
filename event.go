package pokerface

import (
	"fmt"

	"github.com/cfsghost/pokerface/task"
)

type GameEvent int32

const (

	// Initialization
	GameEvent_Started GameEvent = iota
	GameEvent_Initialized
	GameEvent_Prepared
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

var GameEventSymbols = map[GameEvent]string{
	GameEvent_Started:             "Started",
	GameEvent_Initialized:         "Initialized",
	GameEvent_Prepared:            "Prepared",
	GameEvent_AnteRequested:       "AnteRequested",
	GameEvent_AnteReceived:        "AnteReceived",
	GameEvent_Dealt:               "Dealt",
	GameEvent_WagerRequested:      "WagerRequested",
	GameEvent_PreflopRoundEntered: "PreflopRoundEntered",
	GameEvent_FlopRoundEntered:    "FlopRoundEntered",
	GameEvent_TurnRoundEntered:    "TurnRoundEntered",
	GameEvent_RiverRoundEntered:   "RiverRoundEntered",
	GameEvent_RoundInitialized:    "RoundInitialized",
	GameEvent_RoundPrepared:       "RoundPrepared",
	GameEvent_RoundClosed:         "RoundClosed",
	GameEvent_PlayerDidAction:     "PlayerDidAction",
	GameEvent_GameCompleted:       "GameCompleted",
	GameEvent_SettlementRequested: "SettlementRequested",
	GameEvent_SettlementCompleted: "SettlementCompleted",
	GameEvent_GameClosed:          "GameClosed",
}

var GameEventBySymbol = map[string]GameEvent{
	"Started":             GameEvent_Started,
	"Initialized":         GameEvent_Initialized,
	"Prepared":            GameEvent_Prepared,
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

type EventPayload struct {
	Task task.TaskManager `json:"task"`
}

type RoundInitializedEventRuntime struct {
	Dealer int64 `json:"dealer"`
	SB     int64 `json:"sb"`
	BB     int64 `json:"bb"`
}

func NewEventPayload() *EventPayload {
	return &EventPayload{}
}

func (g *game) triggerEvent(event GameEvent) error {

	switch event {

	case GameEvent_Started:
		fmt.Println("Game has started.")
		return g.onStarted()

	case GameEvent_Initialized:
		fmt.Println("Game has been initialized.")
		return g.onInitialized()

	case GameEvent_Prepared:
		fmt.Println("Game has been prepared.")
		return g.onPrepared()

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

func (g *game) EmitEvent(event GameEvent, payload *EventPayload) error {

	// Update current event
	g.gs.Status.CurrentEvent.Name = GameEventSymbols[event]

	if payload != nil {
		g.gs.Status.CurrentEvent.Payload = payload
	} else {
		// Create a new payload for this event
		g.gs.Status.CurrentEvent.Payload = NewEventPayload()
	}

	g.AssertEventRuntime()

	return g.triggerEvent(event)
}

func (g *game) AssertEventRuntime() {

	/*
		eventType, ok := GameEventBySymbol[g.gs.Status.CurrentEvent.Name]
		if !ok {
			return
		}
		   switch eventType {
		   case GameEvent_Initialized:

		   		g.gs.Status.CurrentEvent.Payload = &RoundInitializedEventRuntime{}
		   	}
	*/
}

func (g *game) GetEvent() *Event {
	return g.gs.Status.CurrentEvent
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

	return g.Initialize()
}

func (g *game) onInitialized() error {
	return g.Prepare()
}

func (g *game) onPrepared() error {

	if g.gs.Meta.Ante > 0 {
		return g.RequestAnte()
	}

	return g.EnterPreflopRound()
}

func (g *game) onAnteRequested() error {
	/*
		wg, err := g.WaitForAllPlayersPaidAnte()
		if err != nil {
			return err
		}

		if wg.IsCompleted() {
			g.wg = nil
			return g.EmitEvent(GameEvent_AnteReceived, nil)
		}
	*/
	// Nothing to do, just waiting for all players paid

	return nil
}

func (g *game) onAnteReceived() error {

	// Update pots
	err := g.updatePots()
	if err != nil {
		return err
	}

	g.ResetAllPlayerStatus()

	return g.EmitEvent(GameEvent_PreflopRoundEntered, nil)
}

func (g *game) onPreflopRoundEntered() error {

	g.gs.Status.Round = "preflop"

	return g.InitializeRound()
}

func (g *game) onRoundInitialized() error {
	return g.PrepareRound()
}

func (g *game) onRoundPrepared() error {

	return nil

	event := g.GetEvent()

	event.Payload.Task.Execute()

	if !event.Payload.Task.IsCompleted() {

		// Keep going to wait for ready
		task := event.Payload.Task.GetAvailableTask()

		// Update allowed actions of players based on task state
		players := task.GetPayload().(map[int]bool)
		for idx, isReady := range players {

			if isReady {
				g.Player(idx).AllowActions([]string{})
				continue
			}

			// Not ready so we are waiting for this player
			g.Player(idx).AllowActions([]string{
				"ready",
			})
		}

		fmt.Println("Waiting for ready")

		return nil
	}

	g.ResetAllPlayerAllowedActions()

	return nil
}
