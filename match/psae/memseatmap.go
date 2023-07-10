package psae

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type MemorySeatMap struct {
	tables       map[string]*TableState
	ordered      *list.List
	mu           sync.RWMutex
	totalPlayers int64
}

func NewMemorySeatMap() *MemorySeatMap {
	return &MemorySeatMap{
		tables:       make(map[string]*TableState),
		ordered:      list.New(),
		totalPlayers: 0,
	}
}

func (msm *MemorySeatMap) CreateTable(ts *TableState) error {

	msm.mu.Lock()
	defer msm.mu.Unlock()

	newts := &TableState{
		ID:             ts.ID,
		Players:        ts.Players,
		AvailableSeats: ts.AvailableSeats,
		TotalSeats:     ts.TotalSeats,
		Status:         TableStatus_Ready,
		Statistics: &TableStatistics{
			NoChanges: 0,
		},
	}

	msm.tables[newts.ID] = newts

	if len(newts.Players) > 0 {
		atomic.AddInt64(&msm.totalPlayers, int64(len(newts.Players)))
	}

	// For ordering
	for e := msm.ordered.Front(); e != nil; e = e.Next() {
		et := e.Value.(*TableState)
		if len(et.Players) >= len(ts.Players) {
			msm.ordered.InsertBefore(newts, e)
			return nil
		}
	}

	msm.ordered.PushBack(newts)

	return nil
}

func (msm *MemorySeatMap) DestroyTable(tid string) error {

	msm.mu.Lock()
	defer msm.mu.Unlock()

	ts, ok := msm.tables[tid]
	if !ok {
		return nil
	}

	if len(ts.Players) > 0 {
		atomic.AddInt64(&msm.totalPlayers, -int64(len(ts.Players)))
	}

	delete(msm.tables, tid)

	// For ordering
	for e := msm.ordered.Front(); e != nil; e = e.Next() {
		et := e.Value.(*TableState)
		if et.ID == tid {
			msm.ordered.Remove(e)
			break
		}
	}

	return nil
}

func (msm *MemorySeatMap) UpdateTable(newts *TableState) (*TableState, error) {

	msm.mu.Lock()
	defer msm.mu.Unlock()

	ts, ok := msm.tables[newts.ID]
	if !ok {
		return nil, nil
	}

	// Update states
	if ts.AvailableSeats != newts.AvailableSeats {
		ts.AvailableSeats = newts.AvailableSeats
	}

	if ts.TotalSeats != newts.TotalSeats {
		ts.TotalSeats = newts.TotalSeats
	}

	if ts.Status != newts.Status {
		ts.Status = newts.Status
	}

	pGameIsChanged := false
	if ts.LastGameID != newts.LastGameID {
		pGameIsChanged = true
		//fmt.Printf("last game: %s, new game: %s\n", ts.LastGameID, newts.LastGameID)
		ts.LastGameID = newts.LastGameID
	}

	pNew := 0
	pLeft := 0
	pChanged := false

	for pid, _ := range ts.Players {
		if _, ok := newts.Players[pid]; !ok {

			// player has left
			pLeft++

			pChanged = true
		}
	}

	for pid, _ := range newts.Players {
		if _, ok := ts.Players[pid]; !ok {

			// New player
			pNew++

			pChanged = true
		}
	}

	ts.Players = newts.Players

	atomic.AddInt64(&msm.totalPlayers, int64(pNew-pLeft))

	// Table was changed, it needs to be re-ordered
	if pChanged {

		ts.Statistics.NoChanges = 0

		// Find element
		var tsEle *list.Element
		for e := msm.ordered.Front(); e != nil; e = e.Next() {
			et := e.Value.(*TableState)
			if et.ID == ts.ID {
				tsEle = e
				break
			}
		}

		// Move
		for e := msm.ordered.Front(); e != nil; e = e.Next() {
			et := e.Value.(*TableState)
			if tsEle.Value.(*TableState).ID != et.ID && len(et.Players) >= len(ts.Players) {
				//fmt.Printf("MOVE table %s to new place before table %s\n", tsEle.Value.(*TableState).ID, et.ID)
				msm.ordered.MoveBefore(tsEle, e)
				break
			}
		}

		return ts, nil
	}

	// Count once for the same game
	if pGameIsChanged {
		ts.Statistics.NoChanges++
	}

	return ts, nil
}

func (msm *MemorySeatMap) GetTableState(tid string) (*TableState, error) {

	msm.mu.RLock()
	defer msm.mu.RUnlock()

	ts, ok := msm.tables[tid]
	if !ok {
		return nil, nil
	}

	return ts.Clone(), nil
}

func (msm *MemorySeatMap) FindAvailableTable(condition *TableCondition) (*TableState, error) {

	msm.mu.RLock()
	defer msm.mu.RUnlock()

	if condition.HighestNumberOfPlayers {

		for e := msm.ordered.Back(); e != nil; e = e.Prev() {
			et := e.Value.(*TableState)

			if et.Status != TableStatus_Ready {
				continue
			}

			if condition.MinAvailableSeats == -1 && et.TotalSeats > len(et.Players) {
				return et, nil
			} else if et.AvailableSeats >= condition.MinAvailableSeats {
				return et, nil
			}
		}
	} else {
		for e := msm.ordered.Front(); e != nil; e = e.Next() {
			et := e.Value.(*TableState)

			if et.Status != TableStatus_Ready {
				continue
			}

			//fmt.Printf("FindAvailableTable: %s, player_count=%d, avail=%d\n", et.ID, len(et.Players), et.AvailableSeats)
			if condition.MinAvailableSeats == -1 && et.TotalSeats > len(et.Players) {
				return et, nil
			} else if et.AvailableSeats >= condition.MinAvailableSeats {
				return et, nil
			}
		}
	}

	return nil, nil
}

func (msm *MemorySeatMap) GetTableCount() (int, error) {
	return len(msm.tables), nil
}

func (msm *MemorySeatMap) GetTotalPlayers() (int64, error) {
	return msm.totalPlayers, nil
}

func (msm *MemorySeatMap) GetAllTables() ([]*TableState, error) {

	tables := make([]*TableState, 0, len(msm.tables))

	for _, t := range msm.tables {
		tables = append(tables, t)
	}

	return tables, nil
}
