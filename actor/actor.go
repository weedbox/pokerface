package actor

import pokertable "github.com/weedbox/pokertable/model"

type Actor interface {
	SetRunner(r Runner) error
	SetAdapter(tc Adapter) error
	GetTable() Adapter
	GetRunner() Runner
	UpdateTableState(t *pokertable.Table) error
	GetSeatIndex(playerID string) int
}

type actor struct {
	runner    Runner
	table     Adapter
	tableInfo *pokertable.Table
}

func NewActor() Actor {
	return &actor{}
}

func (a *actor) SetRunner(r Runner) error {
	r.SetActor(a)
	a.runner = r
	return nil
}

func (a *actor) SetAdapter(tc Adapter) error {
	tc.SetActor(a)
	a.table = tc
	return nil
}

func (a *actor) GetTable() Adapter {
	return a.table
}

func (a *actor) GetRunner() Runner {
	return a.runner
}

func (a *actor) UpdateTableState(tableInfo *pokertable.Table) error {
	a.tableInfo = tableInfo
	return a.runner.UpdateTableState(tableInfo)
}

// TODO: this function should be moved to TableStatus struct
func (a *actor) GetSeatIndex(playerID string) int {

	// Find seat index about me
	for _, ps := range a.tableInfo.State.PlayerStates {
		if ps.PlayerID == playerID {
			return ps.SeatIndex
		}
	}

	return -1
}
