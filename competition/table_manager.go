package competition

import (
	"sync"
	"sync/atomic"

	"github.com/weedbox/pokerface/match"
	"github.com/weedbox/pokerface/table"
)

type TableManager interface {
	Initialize() error
	CreateTable() (*table.State, error)
	ReleaseTable(tableID string) error
	ActivateTable(tableID string) error
	GetTables() []*table.State
	GetTableState(tableID string) *table.State
	GetTableCount() int64
	SetJoinable(joinable bool) error
	UpdateTableState(ts *table.State) error
	ReserveSeat(tableID string, seatID int, p *PlayerInfo) (int, error)
	OnTableStateUpdated(fn func(ts *table.State))
	OnSeatChanged(fn func(ts *table.State, sc *match.SeatChanges))
}

type tableManager struct {
	options             *Options
	b                   TableBackend
	tables              map[string]*table.State
	count               int64
	mu                  sync.RWMutex
	onTableStateUpdated func(ts *table.State)
	onSeatChanged       func(ts *table.State, sc *match.SeatChanges)
}

func NewTableManager(options *Options, b TableBackend) TableManager {
	return &tableManager{
		options:             options,
		b:                   b,
		tables:              make(map[string]*table.State),
		onTableStateUpdated: func(ts *table.State) {},
		onSeatChanged:       func(ts *table.State, sc *match.SeatChanges) {},
	}
}

func (tm *tableManager) Initialize() error {

	return nil
	/*
	   // Only one static table
	   if tm.options.MaxTables == 1 {

	   		// Create table immediately
	   		ts, err := tm.CreateTable()
	   		if err != nil {
	   			return err
	   		}

	   		// Then activate
	   		return tm.ActivateTable(ts.ID)
	   	}

	   return nil
	*/
}

func (tm *tableManager) CreateTable() (*table.State, error) {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	opts := table.NewOptions()
	opts.GameType = tm.options.GameType
	opts.MaxSeats = tm.options.Table.MaxSeats
	opts.Ante = tm.options.Table.Ante
	opts.Blind.Dealer = tm.options.Table.Blind.Dealer
	opts.Blind.SB = tm.options.Table.Blind.SB
	opts.Blind.BB = tm.options.Table.Blind.BB
	opts.EliminateMode = "leave"

	ts, err := tm.b.CreateTable(opts)
	if err != nil {
		return nil, err
	}

	tm.tables[ts.ID] = ts

	atomic.AddInt64(&tm.count, 1)

	return ts, nil
}

func (tm *tableManager) ReleaseTable(tableID string) error {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	err := tm.b.ReleaseTable(tableID)
	if err != nil {
		return err
	}

	delete(tm.tables, tableID)

	atomic.AddInt64(&tm.count, -1)

	return nil
}

func (tm *tableManager) ActivateTable(tableID string) error {
	return tm.b.ActivateTable(tableID)
}

func (tm *tableManager) UpdateTableState(ts *table.State) error {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	oldts, ok := tm.tables[ts.ID]
	if !ok {
		return ErrNotFoundTable
	}

	tm.tables[ts.ID] = ts

	sc := tm.GetSeatChanges(oldts, ts)
	if sc != nil {
		go tm.onSeatChanged(ts, sc)
	}

	go tm.onTableStateUpdated(ts)

	return nil
}

func (tm *tableManager) GetSeatChanges(oldts *table.State, newts *table.State) *match.SeatChanges {

	//fmt.Println(newts.Status)

	//newts.PrintState()
	if newts.GameState == nil {
		return nil
	}

	/*
		fmt.Printf("TableUpdated (table=%s, players=%d)\n", ts.ID, len(ts.Players))
		for _, p := range ts.Players {
			fmt.Printf("  origin table (table=%s, seat=%d, id=%s)\n", ts.ID, p.SeatID, p.ID)
		}

		for _, p := range newts.Players {
			fmt.Printf("  new table (table=%s, seat=%d, id=%s)\n", newts.ID, p.SeatID, p.ID)
		}
	*/
	// Preparing seat changes
	sc := match.NewSeatChanges()
	for _, p := range oldts.Players {

		found := false
		for _, np := range newts.Players {
			if p.ID == np.ID {
				found = true
				break
			}
		}

		// Unable to find the player indicates that the player has left the seat
		if !found {
			// Remove players
			sc.Seats[p.SeatID] = "left"
		}
	}

	// Update dealer, sb and bb
	for _, np := range newts.Players {
		if np.CheckPosition("dealer") {
			sc.Dealer = np.SeatID
		}

		if np.CheckPosition("sb") {
			sc.SB = np.SeatID
		}

		if np.CheckPosition("bb") {
			sc.BB = np.SeatID
		}
	}
	/*
		for seatID, _ := range sc.Seats {
			fmt.Printf("    remove player (table=%s, seat=%d)\n", newts.ID, seatID)
		}
	*/

	return sc
}

func (tm *tableManager) GetTableState(tableID string) *table.State {

	tm.mu.RLock()
	defer tm.mu.RUnlock()

	ts, ok := tm.tables[tableID]
	if !ok {
		return nil
	}

	return ts
}

func (tm *tableManager) GetTables() []*table.State {

	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tables := make([]*table.State, 0)

	for _, ts := range tm.tables {
		tables = append(tables, ts)
	}

	return tables
}

func (tm *tableManager) GetTableCount() int64 {
	return atomic.LoadInt64(&tm.count)
}

func (tm *tableManager) SetJoinable(joinable bool) error {

	for _, t := range tm.tables {
		tm.b.SetJoinable(t.ID, joinable)
	}

	return nil
}

func (tm *tableManager) ReserveSeat(tableID string, seatID int, p *PlayerInfo) (int, error) {
	return tm.b.ReserveSeat(tableID, seatID, p)
}

func (tm *tableManager) OnTableStateUpdated(fn func(ts *table.State)) {
	tm.onTableStateUpdated = fn
}

func (tm *tableManager) OnSeatChanged(fn func(ts *table.State, sc *match.SeatChanges)) {
	tm.onSeatChanged = fn
}
