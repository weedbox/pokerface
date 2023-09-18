package match

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type TableCondition struct {
	HighestNumberOfPlayers bool `json:"highest_number_of_players"`
	MinAvailableSeats      int  `json:"min_available_seats"`
}

type TableMap interface {
	GetTables() (map[string]*Table, error)
	CreateTable(players []string) (*Table, error)
	BreakTable(tableID string) error
	GetTable(tableID string) (*Table, error)
	ApplySeatChanges(tableID string, sc *SeatChanges) error
	FindAvailableTable(condition *TableCondition) (*Table, error)
	DispatchPlayer(condition *TableCondition, playerID string) error
	Count() int64

	OnTableBroken(func(*Table))
	OnPlayerLeft(func(table *Table, seatID int, playerID string))
	OnPlayerJoined(func(table *Table, seatID int, playerID string))
}

type tableMap struct {
	m          Match
	tables     map[string]*Table
	ordered    *list.List
	mu         sync.RWMutex
	tableCount int64

	onTableBroken  func(*Table)
	onPlayerJoined func(*Table, int, string)
	onPlayerLeft   func(*Table, int, string)
}

func NewTableMap(m Match) TableMap {
	return &tableMap{
		m:              m,
		tables:         make(map[string]*Table),
		ordered:        list.New(),
		onTableBroken:  func(*Table) {},
		onPlayerJoined: func(*Table, int, string) {},
		onPlayerLeft:   func(*Table, int, string) {},
	}
}

func (tm *tableMap) isTableAvailable(condition *TableCondition, table *Table) bool {
	/*
		fmt.Printf("isAvailableTable: minAvailSeats=%d, table_id=%s, status=%d, player_count=%d, seats=%d, avail_seats=%d\n",
			condition.MinAvailableSeats,
			table.ID(),
			table.status,
			table.GetPlayerCount(),
			table.GetSeatCount(),
			table.GetAvailableSeatCount(),
		)
	*/
	if table.GetStatus() != TableStatus_Ready {
		return false
	}
	/*
		for _, s := range table.sm.GetSeats() {
			fmt.Println("==", s.ID, s.Player, s.IsActive, s.IsReserved)
		}
	*/
	if condition.MinAvailableSeats == -1 {
		if table.GetSeatCount() > table.GetPlayerCount() {
			// players are allowed to sit at any empty seat
			//fmt.Println("table.GetSeatCount() > table.GetPlayerCount() FOUND", table.ID(), table.GetSeatCount(), table.GetPlayerCount())
			return true
		}
	} else if table.GetAvailableSeatCount() >= condition.MinAvailableSeats {
		// minimum number of available seats is required to allow players to sit
		//fmt.Println("table.GetAvailableSeatCount() FOUND", table.ID(), table.GetAvailableSeatCount())
		return true
	}

	return false
}

func (tm *tableMap) findAvailableTable(condition *TableCondition) (*Table, error) {

	if condition.HighestNumberOfPlayers {

		// Start searching from the table with the highest number of players
		for e := tm.ordered.Back(); e != nil; e = e.Prev() {

			var table *Table = e.Value.(*Table)
			if tm.isTableAvailable(condition, table) {
				//fmt.Printf("Found available table: %s\n", table.ID())
				return table, nil
			}
		}

		return nil, ErrNotFoundAvailableTable
	}

	for e := tm.ordered.Front(); e != nil; e = e.Next() {
		var table *Table = e.Value.(*Table)
		if tm.isTableAvailable(condition, table) {
			return table, nil
		}
	}

	return nil, ErrNotFoundAvailableTable
}

func (tm *tableMap) createTable(players []string) (*Table, error) {

	// Allocate table
	t, err := tm.m.TableBackend().Allocate(tm.m.Options().MaxSeats)
	if err != nil {
		return nil, err
	}

	t.SetStatus(TableStatus_Busy)

	// Initializing event handlers
	t.OnPlayerJoined(func(playerID string, seatID int) {
		err := tm.m.TableBackend().Reserve(t.ID(), seatID, playerID)
		if err != nil {
			//fmt.Println("[match/table_map]", playerID, seatID, err)
			return
		}

		tm.onPlayerJoined(t, seatID, playerID)
	})

	t.OnPlayerLeft(func(playerID string, seatID int) {
		tm.onPlayerLeft(t, seatID, playerID)
	})

	// Joins
	for _, playerID := range players {
		t.Join(-1, playerID)
	}

	tm.tables[t.id] = t

	atomic.AddInt64(&tm.tableCount, 1)

	// Update the ordered list
	for e := tm.ordered.Front(); e != nil; e = e.Next() {
		et := e.Value.(*Table)
		if et.GetPlayerCount() >= t.GetPlayerCount() {
			tm.ordered.InsertBefore(t, e)
			return t, nil
		}
	}

	tm.ordered.PushBack(t)

	t.SetStatus(TableStatus_Ready)

	return t, nil
}

func (tm *tableMap) breakTable(tableID string) error {

	t, ok := tm.tables[tableID]
	if !ok {
		return nil
	}

	//fmt.Printf("Break table: %s\n", tableID)

	t.SetStatus(TableStatus_Broken)

	// Attempt to close the target table and release the players
	err := tm.m.TableBackend().Release(tableID)
	if err != nil {
		// Do nothing if failed to release table
		t.SetStatus(TableStatus_Busy)
		return err
	}

	// Delete
	delete(tm.tables, tableID)

	// Update ordered list
	for e := tm.ordered.Front(); e != nil; e = e.Next() {
		et := e.Value.(*Table)
		if et.id == tableID {
			tm.ordered.Remove(e)
			break
		}
	}

	atomic.AddInt64(&tm.tableCount, -1)

	// Release table
	t.Release()

	// Dismiss
	err = tm.m.Runner().DismissTable(tm.m, t)
	if err != nil {
		return err
	}

	tm.onTableBroken(t)

	return nil
}

func (tm *tableMap) getTable(tableID string) (*Table, error) {

	t, ok := tm.tables[tableID]
	if !ok {
		return nil, ErrNotFoundTable
	}

	return t, nil
}

func (tm *tableMap) GetTables() (map[string]*Table, error) {
	return tm.tables, nil
}

func (tm *tableMap) CreateTable(players []string) (*Table, error) {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, err := tm.createTable(players)

	// Activate table
	err = tm.m.TableBackend().Activate(t.id)
	if err != nil {
		return nil, err
	}

	t.SetStatus(TableStatus_Ready)

	return t, err
}

func (tm *tableMap) BreakTable(tableID string) error {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	return tm.breakTable(tableID)
}

func (tm *tableMap) GetTable(tableID string) (*Table, error) {

	tm.mu.RLock()
	defer tm.mu.RUnlock()

	return tm.getTable(tableID)
}

func (tm *tableMap) ApplySeatChanges(tableID string, sc *SeatChanges) error {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, err := tm.getTable(tableID)
	if err != nil {
		return err
	}

	t.SetStatus(TableStatus_Busy)

	err = t.ApplySeatChanges(sc)
	if err != nil {
		return err
	}

	// Does the table need to continue to exist?
	if tm.m.Runner().ShouldBeSplit(tm.m, t) {

		// Break it
		err := tm.breakTable(t.ID())
		if err != nil {
			return err
		}

		return nil
	}

	t.SetStatus(TableStatus_Ready)

	return nil
}

func (tm *tableMap) Count() int64 {

	return atomic.LoadInt64(&tm.tableCount)
}

func (tm *tableMap) FindAvailableTable(condition *TableCondition) (*Table, error) {

	tm.mu.RLock()
	defer tm.mu.RUnlock()

	return tm.findAvailableTable(condition)
}

func (tm *tableMap) DispatchPlayer(condition *TableCondition, playerID string) error {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Find the table with the maximum number of players
	table, err := tm.findAvailableTable(condition)
	if err == nil && table != nil {
		// Found a available table can be used
		return table.Join(-1, playerID)
	}

	return err
}

func (tm *tableMap) OnTableBroken(fn func(*Table)) {
	tm.onTableBroken = fn
}

func (tm *tableMap) OnPlayerJoined(fn func(table *Table, seatID int, playerID string)) {
	tm.onPlayerJoined = fn
}

func (tm *tableMap) OnPlayerLeft(fn func(table *Table, seatID int, playerID string)) {
	tm.onPlayerLeft = fn
}
