package actor

import (
	"github.com/weedbox/pokertable"
)

type ObserverRunner struct {
	actor               Actor
	tableInfo           *pokertable.Table
	systemMode          bool
	onTableStateUpdated func(*pokertable.Table)
}

func NewObserverRunner() *ObserverRunner {
	return &ObserverRunner{
		onTableStateUpdated: func(*pokertable.Table) {},
	}
}

func (obr *ObserverRunner) SetActor(a Actor) {
	obr.actor = a
}

func (obr *ObserverRunner) EnabledSystemMode(enabled bool) {
	obr.systemMode = enabled
}

func (obr *ObserverRunner) UpdateTableState(tableInfo *pokertable.Table) error {

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

func (obr *ObserverRunner) OnTableStateUpdated(fn func(*pokertable.Table)) error {
	obr.onTableStateUpdated = fn
	return nil
}
