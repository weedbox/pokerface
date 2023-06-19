package actor

import (
	"github.com/weedbox/pokertable"
)

type observerRunner struct {
	actor               Actor
	tableInfo           *pokertable.Table
	systemMode          bool
	onTableStateUpdated func(*pokertable.Table)
}

func NewObserverRunner() *observerRunner {
	return &observerRunner{
		onTableStateUpdated: func(*pokertable.Table) {},
	}
}

func (obr *observerRunner) SetActor(a Actor) {
	obr.actor = a
}

func (obr *observerRunner) EnabledSystemMode(enabled bool) {
	obr.systemMode = enabled
}

func (obr *observerRunner) UpdateTableState(tableInfo *pokertable.Table) error {

	obr.tableInfo = tableInfo

	if !obr.systemMode {
		// Filtering private information for observer
		switch tableInfo.State.Status {
		case pokertable.TableStateStatus_TableGamePlaying:
			fallthrough
		case pokertable.TableStateStatus_TableGameSettled:
			tableInfo.State.GameState.AsObserver()
		}
	}

	// Emit event
	obr.onTableStateUpdated(tableInfo)

	return nil
}

func (obr *observerRunner) OnTableStateUpdated(fn func(*pokertable.Table)) error {
	obr.onTableStateUpdated = fn
	return nil
}
