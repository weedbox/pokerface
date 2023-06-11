package actor

import (
	"time"

	"github.com/weedbox/pokerface"
	"github.com/weedbox/pokertable"
	"github.com/weedbox/timebank"
)

type playerRunner struct {
	actor               Actor
	actions             Actions
	playerID            string
	gamePlayerIdx       int
	tableInfo           *pokertable.Table
	timebank            *timebank.TimeBank
	onTableStateUpdated func(*pokertable.Table)
}

func NewPlayerRunner(playerID string) *playerRunner {
	return &playerRunner{
		playerID: playerID,
		timebank: timebank.NewTimeBank(),
	}
}

func (pr *playerRunner) SetActor(a Actor) {
	pr.actor = a
	pr.actions = NewActions(a, pr.playerID)
}

func (pr *playerRunner) UpdateTableState(table *pokertable.Table) error {

	// Update seat index
	pr.gamePlayerIdx = table.GamePlayerIndex(pr.playerID)

	// Filtering private information fpr player
	table.State.GameState.AsPlayer(pr.gamePlayerIdx)

	pr.tableInfo = table

	// Emit event
	pr.onTableStateUpdated(table)

	// Somehow, this player is not in the game.
	// It probably has no chips already.
	if pr.gamePlayerIdx == -1 {
		return nil
	}

	// Game is running right now
	switch pr.tableInfo.State.Status {
	case pokertable.TableStateStatus_TableGameMatchOpen:

		// We have actions allowed by game engine
		player := pr.tableInfo.State.GameState.GetPlayer(pr.gamePlayerIdx)
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

	gs := pr.tableInfo.State.GameState

	// Do pass automatically
	if gs.HasAction(pr.gamePlayerIdx, "pass") {
		return pr.actions.Pass()
	}

	// Setup timebank to wait for player
	thinkingTime := time.Duration(pr.tableInfo.Meta.CompetitionMeta.ActionTimeSecs) * time.Second
	return pr.timebank.NewTask(thinkingTime, func(isCancelled bool) {

		if isCancelled {
			return
		}

		// Do default actions if player has no response
		pr.automate()
	})
}

func (pr *playerRunner) automate() error {

	gs := pr.tableInfo.State.GameState

	// Default actions for automation when player has no response
	if gs.HasAction(pr.gamePlayerIdx, "ready") {
		return pr.actions.Ready()
	} else if gs.HasAction(pr.gamePlayerIdx, "check") {
		return pr.actions.Check()
	} else if gs.HasAction(pr.gamePlayerIdx, "fold") {
		return pr.actions.Fold()
	}

	// Pay for ante and blinds
	switch gs.Status.CurrentEvent.Name {
	case pokerface.GameEventSymbols[pokerface.GameEvent_Prepared]:

		// Ante
		return pr.actions.Pay(gs.Meta.Ante)

	case pokerface.GameEventSymbols[pokerface.GameEvent_RoundInitialized]:

		// blinds
		if gs.HasPosition(pr.gamePlayerIdx, "sb") {
			return pr.actions.Pay(gs.Meta.Blind.SB)
		} else if gs.HasPosition(pr.gamePlayerIdx, "bb") {
			return pr.actions.Pay(gs.Meta.Blind.BB)
		}

		return pr.actions.Pay(gs.Meta.Blind.Dealer)
	}

	return nil
}
