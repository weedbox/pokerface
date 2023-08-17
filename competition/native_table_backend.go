package competition

import (
	"sync"

	"github.com/weedbox/pokerface/table"
)

type NativeTableBackend struct {
	geb            table.Backend
	tables         map[string]table.Table
	mu             sync.RWMutex
	onTableUpdated func(ts *table.State)
}

func NewNativeTableBackend(geb table.Backend) TableBackend {
	return &NativeTableBackend{
		geb:            geb,
		tables:         make(map[string]table.Table),
		onTableUpdated: func(ts *table.State) {},
	}
}

func (ntb *NativeTableBackend) CreateTable(opts *table.Options) (*table.State, error) {

	ntb.mu.Lock()
	defer ntb.mu.Unlock()

	t := table.NewTable(opts, table.WithBackend(ntb.geb))

	t.OnStateUpdated(func(ts *table.State) {
		ntb.onTableUpdated(ts)
	})

	ntb.tables[t.GetState().ID] = t

	return t.GetState(), nil
}

func (ntb *NativeTableBackend) ActivateTable(tableID string) error {

	ntb.mu.RLock()
	defer ntb.mu.RUnlock()

	t, ok := ntb.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	return t.Start()
}

func (ntb *NativeTableBackend) ReleaseTable(tableID string) error {

	ntb.mu.Lock()
	defer ntb.mu.Unlock()

	t, ok := ntb.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	delete(ntb.tables, tableID)

	return t.Close()
}

func (ntb *NativeTableBackend) SetJoinable(tableID string, isJoinable bool) error {

	ntb.mu.Lock()
	defer ntb.mu.Unlock()

	t, ok := ntb.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	t.SetJoinable(isJoinable)

	return nil
}

func (ntb *NativeTableBackend) ReserveSeat(tableID string, seatID int, p *PlayerInfo) (int, error) {

	ntb.mu.RLock()
	defer ntb.mu.RUnlock()

	t, ok := ntb.tables[tableID]
	if !ok {
		return -1, ErrNotFoundTable
	}

	seatID, err := t.Join(seatID, &table.PlayerInfo{
		ID:       p.ID,
		Bankroll: p.Bankroll,
	})

	return seatID, err
}

func (ntb *NativeTableBackend) OnTableUpdated(fn func(ts *table.State)) {
	ntb.onTableUpdated = fn
}

// For testing
func (ntb *NativeTableBackend) GetTable(tableID string) table.Table {

	ntb.mu.RLock()
	defer ntb.mu.RUnlock()

	t, ok := ntb.tables[tableID]
	if !ok {
		return nil
	}

	return t
}
