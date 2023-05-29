package pokerface

import (
	"github.com/cfsghost/pokerface/task"
)

type GameEvent int32

const (

	// Initialization
	GameEvent_Started GameEvent = iota
	GameEvent_Initialized
	GameEvent_Prepared
	GameEvent_AnteRequested

	// States
	GameEvent_PlayerActionRequested

	// Rounds
	GameEvent_PreflopRoundEntered
	GameEvent_FlopRoundEntered
	GameEvent_TurnRoundEntered
	GameEvent_RiverRoundEntered
	GameEvent_RoundInitialized
	GameEvent_RoundPrepared
	GameEvent_RoundClosed

	// Result
	GameEvent_GameCompleted
	GameEvent_SettlementRequested
	GameEvent_SettlementCompleted
	GameEvent_GameClosed
)

var GameEventSymbols = map[GameEvent]string{
	GameEvent_Started:               "Started",
	GameEvent_Initialized:           "Initialized",
	GameEvent_Prepared:              "Prepared",
	GameEvent_AnteRequested:         "AnteRequested",
	GameEvent_PlayerActionRequested: "PlayerActionRequested",
	GameEvent_PreflopRoundEntered:   "PreflopRoundEntered",
	GameEvent_FlopRoundEntered:      "FlopRoundEntered",
	GameEvent_TurnRoundEntered:      "TurnRoundEntered",
	GameEvent_RiverRoundEntered:     "RiverRoundEntered",
	GameEvent_RoundInitialized:      "RoundInitialized",
	GameEvent_RoundPrepared:         "RoundPrepared",
	GameEvent_RoundClosed:           "RoundClosed",
	GameEvent_GameCompleted:         "GameCompleted",
	GameEvent_SettlementRequested:   "SettlementRequested",
	GameEvent_SettlementCompleted:   "SettlementCompleted",
	GameEvent_GameClosed:            "GameClosed",
}

var GameEventBySymbol = map[string]GameEvent{
	"Started":               GameEvent_Started,
	"Initialized":           GameEvent_Initialized,
	"Prepared":              GameEvent_Prepared,
	"AnteRequested":         GameEvent_AnteRequested,
	"PlayerActionRequested": GameEvent_PlayerActionRequested,
	"PreflopRoundEntered":   GameEvent_PreflopRoundEntered,
	"FlopRoundEntered":      GameEvent_FlopRoundEntered,
	"TurnRoundEntered":      GameEvent_TurnRoundEntered,
	"RiverRoundEntered":     GameEvent_RiverRoundEntered,
	"RoundInitialized":      GameEvent_RoundInitialized,
	"RoundPrepared":         GameEvent_RoundPrepared,
	"RoundClosed":           GameEvent_RoundClosed,
	"GameCompleted":         GameEvent_GameCompleted,
	"SettlementRequested":   GameEvent_SettlementRequested,
	"SettlementCompleted":   GameEvent_SettlementCompleted,
	"GameClosed":            GameEvent_GameClosed,
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
		//fmt.Println("Game has started.")
		return g.onStarted()

	case GameEvent_Initialized:
		//fmt.Println("Game has been initialized.")
		return g.onInitialized()

	case GameEvent_Prepared:
		//fmt.Println("Game has been prepared.")
		return g.onPrepared()

	case GameEvent_AnteRequested:
		//fmt.Println("Ante has been requested.")
		return g.onAnteRequested()

	case GameEvent_PlayerActionRequested:
		//fmt.Println("Player action has been requested.")
		return g.onPlayerActionRequested()

	case GameEvent_PreflopRoundEntered:
		//fmt.Println("Entered Preflop round.")
		return g.onPreflopRoundEntered()

	case GameEvent_FlopRoundEntered:
		//fmt.Println("Entered Flop round.")
		return g.onFlopRoundEntered()

	case GameEvent_TurnRoundEntered:
		//fmt.Println("Entered Turn round.")
		return g.onTurnRoundEntered()

	case GameEvent_RiverRoundEntered:
		//fmt.Println("Entered River round.")
		return g.onRiverRoundEntered()

	case GameEvent_RoundInitialized:
		//fmt.Println("Current round has initialized.")
		return g.onRoundInitialized()

	case GameEvent_RoundPrepared:
		//fmt.Println("Current round has been prepared.")
		return g.onRoundPrepared()

	case GameEvent_RoundClosed:
		//fmt.Println("Current round has closed.")
		return g.onRoundClosed()

	case GameEvent_GameCompleted:
		//fmt.Println("Game has been completed.")
		return g.onGameCompleted()

	case GameEvent_SettlementRequested:
		//fmt.Println("Settlement has been requested.")
		return g.onSettlementRequested()

	case GameEvent_SettlementCompleted:
		//fmt.Println("Settlement has been completed.")
		return g.onSettlementCompleted()

	case GameEvent_GameClosed:
		//fmt.Println("Game has closed.")
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

	return g.triggerEvent(event)
}

func (g *game) GetEvent() *Event {
	return g.gs.Status.CurrentEvent
}

func (g *game) onStarted() error {
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

	// Update pots
	err := g.updatePots()
	if err != nil {
		return err
	}

	g.ResetAllPlayerStatus()

	return g.EmitEvent(GameEvent_PreflopRoundEntered, nil)
}

func (g *game) onPlayerActionRequested() error {
	return g.RequestPlayerAction()
}

func (g *game) onRoundInitialized() error {
	return g.PrepareRound()
}

func (g *game) onRoundPrepared() error {
	return g.EmitEvent(GameEvent_PlayerActionRequested, nil)
}

func (g *game) onRoundClosed() error {

	// Update pots
	err := g.updatePots()
	if err != nil {
		return err
	}

	return nil
}

func (g *game) onPreflopRoundEntered() error {

	g.gs.Status.Round = "preflop"

	return g.InitializeRound()
}

func (g *game) onFlopRoundEntered() error {

	g.gs.Status.Round = "flop"

	return g.InitializeRound()
}

func (g *game) onTurnRoundEntered() error {

	g.gs.Status.Round = "turn"

	return g.InitializeRound()
}

func (g *game) onRiverRoundEntered() error {

	g.gs.Status.Round = "river"

	return g.InitializeRound()
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
