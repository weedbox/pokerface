package actor

import (
	"time"

	"github.com/weedbox/pokerface"
	"github.com/weedbox/pokertable"
	"github.com/weedbox/timebank"
)

type PlayerStatus int32

const (
	PlayerStatus_Running PlayerStatus = iota
	PlayerStatus_Idle
	PlayerStatus_Suspend
)

type playerRunner struct {
	actor               Actor
	actions             Actions
	playerID            string
	gamePlayerIdx       int
	tableInfo           *pokertable.Table
	timebank            *timebank.TimeBank
	onTableStateUpdated func(*pokertable.Table)

	// status
	status           PlayerStatus
	idleCount        int
	suspendThreshold int
}

func NewPlayerRunner(playerID string) *playerRunner {
	return &playerRunner{
		playerID:            playerID,
		timebank:            timebank.NewTimeBank(),
		status:              PlayerStatus_Running,
		suspendThreshold:    2,
		onTableStateUpdated: func(*pokertable.Table) {},
	}
}

func (pr *playerRunner) SetActor(a Actor) {
	pr.actor = a
	pr.actions = NewActions(a, pr.playerID)
}

func (pr *playerRunner) UpdateTableState(table *pokertable.Table) error {

	// Update seat index
	pr.gamePlayerIdx = table.GamePlayerIndex(pr.playerID)

	pr.tableInfo = table

	// Emit event
	pr.onTableStateUpdated(table)

	// Game is running right now
	switch pr.tableInfo.State.Status {
	case pokertable.TableStateStatus_TableGamePlaying:

		// Somehow, this player is not in the game.
		// It probably has no chips already.
		if pr.gamePlayerIdx == -1 {
			return nil
		}

		// Filtering private information fpr player
		table.State.GameState.AsPlayer(pr.gamePlayerIdx)

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

	// Player is suspended
	if pr.status == PlayerStatus_Suspend {
		return pr.automate()
	}

	// Setup timebank to wait for player
	thinkingTime := time.Duration(pr.tableInfo.Meta.CompetitionMeta.ActionTime) * time.Second
	return pr.timebank.NewTask(thinkingTime, func(isCancelled bool) {

		if isCancelled {
			return
		}

		// Stay idle already
		if pr.status == PlayerStatus_Idle {
			pr.Idle()
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

func (pr *playerRunner) SetSuspendThreshold(count int) {
	pr.suspendThreshold = count
}

func (pr *playerRunner) Resume() error {

	if pr.status == PlayerStatus_Running {
		return nil
	}

	pr.status = PlayerStatus_Running
	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Idle() error {
	if pr.status != PlayerStatus_Idle {
		pr.status = PlayerStatus_Idle
		pr.idleCount = 0
	} else {
		pr.idleCount++
	}

	if pr.idleCount == pr.suspendThreshold {
		return pr.Suspend()
	}

	return nil
}

func (pr *playerRunner) Suspend() error {
	pr.status = PlayerStatus_Suspend
	return nil
}

func (pr *playerRunner) Pass() error {

	err := pr.actions.Pass()
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Ready() error {

	err := pr.actions.Ready()
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Pay(chips int64) error {

	err := pr.actions.Pay(chips)
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Check() error {

	err := pr.actions.Check()
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Bet(chips int64) error {

	err := pr.actions.Bet(chips)
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Call() error {

	err := pr.actions.Call()
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Fold() error {
	pr.Resume()
	return pr.actions.Fold()
}

func (pr *playerRunner) Allin() error {

	err := pr.actions.Allin()
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}

func (pr *playerRunner) Raise(chipLevel int64) error {

	err := pr.actions.Raise(chipLevel)
	if err != nil {
		return err
	}

	pr.idleCount = 0

	return nil
}
