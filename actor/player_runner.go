package actor

import (
	"github.com/weedbox/pokerface"
	pokertable "github.com/weedbox/pokertable/model"
)

type playerRunner struct {
	actor               Actor
	actions             Actions
	playerID            string
	seatIndex           int
	tableInfo           *pokertable.Table
	onTableStateUpdated func(*pokertable.Table)
}

func NewPlayerRunner(playerID string) *playerRunner {
	return &playerRunner{
		playerID: playerID,
	}
}

func (pr *playerRunner) SetActor(a Actor) {
	pr.actor = a
	pr.actions = NewActions(a, pr.playerID)
}

func (pr *playerRunner) UpdateTableState(table *pokertable.Table) error {

	// Update seat index
	pr.seatIndex = pr.actor.GetSeatIndex(pr.playerID)

	// Filtering private information fpr player
	table.State.GameState.AsPlayer(pr.seatIndex)

	pr.tableInfo = table

	// Emit event
	pr.onTableStateUpdated(table)

	// Game is running right now
	switch pr.tableInfo.State.Status {
	case pokertable.TableStateStatus_TableGameMatchOpen:

		// We have actions allowed by game engine
		player := pr.tableInfo.State.GameState.GetPlayer(pr.seatIndex)
		if len(player.AllowedActions) > 0 {
			pr.requestMove()
		}
	}

	return nil
}

func (pr *playerRunner) OnTableStateUpdated(fn func(*pokertable.Table)) error {
	pr.onTableStateUpdated = fn
	return nil
}

func (pr *playerRunner) requestMove() error {

	//TODO: Setup timer to wait for player
	return nil

	gs := pr.tableInfo.State.GameState

	// Default actions for automation when player has no response
	if gs.HasAction(pr.seatIndex, "ready") {
		return pr.actions.Ready()
	} else if gs.HasAction(pr.seatIndex, "check") {
		return pr.actions.Check()
	} else if gs.HasAction(pr.seatIndex, "fold") {
		return pr.actions.Fold()
	}

	// Pay for ante and blinds
	switch gs.Status.CurrentEvent.Name {
	case pokerface.GameEventSymbols[pokerface.GameEvent_Prepared]:

		// Ante
		return pr.actions.Pay(gs.Meta.Ante)

	case pokerface.GameEventSymbols[pokerface.GameEvent_RoundInitialized]:

		// blinds
		if gs.HasPosition(pr.seatIndex, "sb") {
			return pr.actions.Pay(gs.Meta.Blind.SB)
		} else if gs.HasPosition(pr.seatIndex, "bb") {
			return pr.actions.Pay(gs.Meta.Blind.BB)
		}

		return pr.actions.Pay(gs.Meta.Blind.Dealer)
	}

	return nil
}
