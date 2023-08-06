package actor

import (
	"sync"

	pokertable "github.com/weedbox/pokertable"
)

type Actor interface {
	SetRunner(r Runner) error
	SetAdapter(tc Adapter) error
	GetTable() Adapter
	GetRunner() Runner
	UpdateTableState(t *pokertable.Table) error
}

type actor struct {
	runner       Runner
	tableAdapter Adapter
	mu           sync.RWMutex
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
	a.tableAdapter = tc
	return nil
}

func (a *actor) GetTable() Adapter {
	return a.tableAdapter
}

func (a *actor) GetRunner() Runner {
	return a.runner
}

func (a *actor) UpdateTableState(tableInfo *pokertable.Table) error {

	a.mu.Lock()
	defer a.mu.Unlock()

	err := a.runner.UpdateTableState(tableInfo)
	if err != nil {
		return err
	}

	return nil
}
