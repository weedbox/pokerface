package actor

import (
	"github.com/weedbox/pokertable"
)

type observerRunner struct {
	actor               Actor
	tableInfo           *pokertable.Table
	onTableStateUpdated func(*pokertable.Table)
	onGameStateUpdated  func(*pokertable.Table)
}

func NewObserverRunner() *observerRunner {
	return &observerRunner{
		onTableStateUpdated: func(*pokertable.Table) {},
		onGameStateUpdated:  func(*pokertable.Table) {},
	}
}

func (obr *observerRunner) SetActor(a Actor) {
	obr.actor = a
}

func (obr *observerRunner) UpdateTableState(tableInfo *pokertable.Table) error {

	obr.tableInfo = tableInfo

	// Filtering private information fobr observer
	if tableInfo.State.Status == pokertable.TableStateStatus_TableGamePlaying {
		tableInfo.State.GameState.AsObserver()

		obr.onGameStateUpdated(tableInfo)
	}

	// Emit event
	obr.onTableStateUpdated(tableInfo)

	return nil
}

func (obr *observerRunner) OnTableStateUpdated(fn func(*pokertable.Table)) error {
	obr.onTableStateUpdated = fn
	return nil
}

func (obr *observerRunner) OnGameStateUpdated(fn func(*pokertable.Table)) error {
	obr.onGameStateUpdated = fn
	return nil
}
