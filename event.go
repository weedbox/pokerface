package pokerface

import (
	"github.com/weedbox/pokerface/task"
)

type GameEvent int32

const (

	// Initialization
	GameEvent_Started GameEvent = iota
	GameEvent_Initialized
	GameEvent_Prepared
	GameEvent_AnteRequested

	// Rounds
	GameEvent_PreflopRoundEntered
	GameEvent_FlopRoundEntered
	GameEvent_TurnRoundEntered
	GameEvent_RiverRoundEntered
	GameEvent_RoundInitialized
	GameEvent_RoundPrepared
	GameEvent_RoundStarted
	GameEvent_RoundClosed

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
	GameEvent_PreflopRoundEntered: "PreflopRoundEntered",
	GameEvent_FlopRoundEntered:    "FlopRoundEntered",
	GameEvent_TurnRoundEntered:    "TurnRoundEntered",
	GameEvent_RiverRoundEntered:   "RiverRoundEntered",
	GameEvent_RoundInitialized:    "RoundInitialized",
	GameEvent_RoundPrepared:       "RoundPrepared",
	GameEvent_RoundStarted:        "RoundStarted",
	GameEvent_RoundClosed:         "RoundClosed",
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
	"PreflopRoundEntered": GameEvent_PreflopRoundEntered,
	"FlopRoundEntered":    GameEvent_FlopRoundEntered,
	"TurnRoundEntered":    GameEvent_TurnRoundEntered,
	"RiverRoundEntered":   GameEvent_RiverRoundEntered,
	"RoundInitialized":    GameEvent_RoundInitialized,
	"RoundPrepared":       GameEvent_RoundPrepared,
	"RoundStarted":        GameEvent_RoundStarted,
	"RoundClosed":         GameEvent_RoundClosed,
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
		return g.onStarted()

	case GameEvent_Initialized:
		return g.onInitialized()

	case GameEvent_Prepared:
		return g.onPrepared()

	case GameEvent_AnteRequested:
		return g.onAnteRequested()

	case GameEvent_RoundStarted:
		return g.onRoundStarted()

	case GameEvent_PreflopRoundEntered:
		return g.onPreflopRoundEntered()

	case GameEvent_FlopRoundEntered:
		return g.onFlopRoundEntered()

	case GameEvent_TurnRoundEntered:
		return g.onTurnRoundEntered()

	case GameEvent_RiverRoundEntered:
		return g.onRiverRoundEntered()

	case GameEvent_RoundInitialized:
		return g.onRoundInitialized()

	case GameEvent_RoundPrepared:
		return g.onRoundPrepared()

	case GameEvent_RoundClosed:
		return g.onRoundClosed()

	case GameEvent_GameCompleted:
		return g.onGameCompleted()

	case GameEvent_SettlementRequested:
		return g.onSettlementRequested()

	case GameEvent_SettlementCompleted:
		return g.onSettlementCompleted()

	case GameEvent_GameClosed:
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

	return g.EnterPreflopRound()
	// return g.EmitEvent(GameEvent_PreflopRoundEntered, nil)
}

func (g *game) onRoundStarted() error {
	return g.RequestPlayerAction()
}

func (g *game) onRoundInitialized() error {
	return g.PrepareRound()
}

func (g *game) onRoundPrepared() error {
	return g.EmitEvent(GameEvent_RoundStarted, nil)
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
	return g.InitializeRound()
}

func (g *game) onFlopRoundEntered() error {
	return g.InitializeRound()
}

func (g *game) onTurnRoundEntered() error {
	return g.InitializeRound()
}

func (g *game) onRiverRoundEntered() error {
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
