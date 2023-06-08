package actor

import (
	pokertable "github.com/weedbox/pokertable/model"
)

type observerRunner struct {
	actor               Actor
	tableInfo           *pokertable.Table
	onTableStateUpdated func(*pokertable.Table)
}

func NewObserverRunner() *observerRunner {
	return &observerRunner{}
}

func (obr *observerRunner) SetActor(a Actor) {
	obr.actor = a
}

func (obr *observerRunner) UpdateTableState(table *pokertable.Table) error {

	// Filtering private information fobr observer
	table.State.GameState.AsObserver()

	obr.tableInfo = table

	// Emit event
	obr.onTableStateUpdated(table)

	return nil
}

func (obr *observerRunner) OnTableStateUpdated(fn func(*pokertable.Table)) error {
	obr.onTableStateUpdated = fn
	return nil
}
