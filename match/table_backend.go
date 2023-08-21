package match

import (
	"sync"
)

type TableBackend interface {
	Allocate(maxSeats int) (*Table, error)
	Release(tableID string) error
	Activate(tableID string) error
	Reserve(tableID string, seatID int, playerID string) error
	GetTable(tableID string) (*Table, error)
	UpdateTable(tableID string, sc *SeatChanges) error
	OnTableUpdated(func(tableID string, sc *SeatChanges))
}

type NativeTableBackend struct {
	tables         map[string]*Table
	mu             sync.RWMutex
	onTableUpdated func(tableID string, sc *SeatChanges)
}

func NewDummyTableBackend() *NativeTableBackend {
	return &NativeTableBackend{
		tables:         make(map[string]*Table),
		onTableUpdated: func(tableID string, sc *SeatChanges) {},
	}
}

// For testing and debugging
func (tb *NativeTableBackend) getTables() map[string]*Table {

	tb.mu.RLock()
	defer tb.mu.RUnlock()

	tables := make(map[string]*Table)

	for id, t := range tb.tables {
		tables[id] = t
	}

	return tables
}

func (tb *NativeTableBackend) Allocate(maxSeats int) (*Table, error) {

	tb.mu.Lock()
	defer tb.mu.Unlock()

	internalState := NewTable(maxSeats)
	tb.tables[internalState.ID()] = internalState

	t := NewTable(maxSeats)
	t.id = internalState.ID()

	//fmt.Printf("Allocate Table (id=%s)\n", internalState.ID())

	return t, nil
}

func (tb *NativeTableBackend) Release(tableID string) error {

	tb.mu.Lock()
	defer tb.mu.Unlock()

	delete(tb.tables, tableID)
	return nil
}

func (tb *NativeTableBackend) Activate(tableID string) error {
	return nil
}

func (tb *NativeTableBackend) Reserve(tableID string, seatID int, playerID string) error {

	tb.mu.RLock()
	defer tb.mu.RUnlock()

	t, ok := tb.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	//fmt.Printf("Reserve Seat (id=%s, player=%s)\n", tableID, playerID)

	return t.Join(seatID, playerID)
}

func (tb *NativeTableBackend) GetTable(tableID string) (*Table, error) {

	tb.mu.RLock()
	defer tb.mu.RUnlock()

	t, ok := tb.tables[tableID]
	if !ok {
		return nil, ErrNotFoundTable
	}

	return t, nil
}

func (tb *NativeTableBackend) UpdateTable(tableID string, sc *SeatChanges) error {
	t, _ := tb.GetTable(tableID)
	t.SetStatus(TableStatus_Busy)
	tb.onTableUpdated(tableID, sc)
	t.SetStatus(TableStatus_Ready)
	return nil
}

func (tb *NativeTableBackend) OnTableUpdated(fn func(tableID string, sc *SeatChanges)) {
	tb.onTableUpdated = fn
}
