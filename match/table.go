package match

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/weedbox/pokerface/seat_manager"
)

type TableStatus int

const (
	TableStatus_Preparing TableStatus = iota
	TableStatus_Ready
	TableStatus_Busy
	TableStatus_Broken
	TableStatus_Suspend
)

type Table struct {
	id        string
	status    TableStatus
	sm        *seat_manager.SeatManager
	noChanges int
	mu        sync.RWMutex

	onPlayerJoined  func(playerID string, seatID int)
	onPlayerLeft    func(playerID string, seatID int)
	onPlayerDrained func(playerID string, seatID int)
}

func NewTable(maxSeats int) *Table {
	return &Table{
		id:              uuid.New().String(),
		status:          TableStatus_Preparing,
		sm:              seat_manager.NewSeatManager(maxSeats),
		onPlayerJoined:  func(string, int) {},
		onPlayerLeft:    func(string, int) {},
		onPlayerDrained: func(string, int) {},
	}
}

func (t *Table) SetID(id string) {
	t.id = id
}

func (t *Table) ID() string {
	return t.id
}

func (t *Table) SeatManager() *seat_manager.SeatManager {
	return t.sm
}

func (t *Table) SetStatus(status TableStatus) error {
	t.status = status
	return nil
}

func (t *Table) GetStatus() TableStatus {
	return t.status
}

func (t *Table) Join(seatID int, playerID string) error {

	t.mu.Lock()
	defer t.mu.Unlock()

	sid, err := t.sm.Join(seatID, playerID)
	if err != nil {
		return err
	}

	t.onPlayerJoined(playerID, sid)

	return nil
}

func (t *Table) Release() error {

	t.mu.Lock()
	defer t.mu.Unlock()

	seats := t.sm.GetSeats()
	for _, s := range seats {
		if s.Player != nil {
			t.onPlayerDrained(s.Player.(string), s.ID)
		}
	}

	return nil
}

func (t *Table) isNewGame(sc *SeatChanges) bool {

	if t.sm.Dealer() == nil || t.sm.Dealer().ID != sc.Dealer ||
		t.sm.SmallBlind() == nil || t.sm.SmallBlind().ID != sc.SB ||
		t.sm.BigBlind() == nil || t.sm.BigBlind().ID != sc.BB {
		return true
	}

	/*
		if t.sm.Dealer().ID != sc.Dealer || t.sm.SmallBlind().ID != sc.SB || t.sm.BigBlind().ID != sc.BB {
			return true
		}
	*/
	return false
}

func (t *Table) ApplySeatChanges(sc *SeatChanges) error {

	t.mu.Lock()
	defer t.mu.Unlock()

	if sc.Dealer > -1 && sc.SB > -1 && sc.BB > -1 {

		if t.isNewGame(sc) {
			t.noChanges++
		} else {
			t.noChanges = 0
		}

		// Update positions
		t.sm.SetDealer(sc.Dealer)
		t.sm.SetSmallBlind(sc.SB)
		t.sm.SetBigBlind(sc.BB)
	}

	// Update seats
	for seatID, state := range sc.Seats {
		if state != "left" {
			continue
		}

		seat := t.sm.GetSeat(seatID)
		if seat.Player == nil {
			fmt.Printf("WARNING seat.Player = nil, seat=%d, table=%s\n", seatID, t.ID())
			continue
		}

		playerID := seat.Player.(string)
		t.sm.Leave(seatID)
		t.onPlayerLeft(playerID, seatID)
	}

	return nil
}

func (t *Table) GetPlayerCount() int {

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.sm.GetPlayerCount()
}

func (t *Table) GetSeatCount() int {

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.sm.GetSeatCount()
}

func (t *Table) GetAvailableSeatCount() int {

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.sm.GetAvailableSeatCount()
}

func (t *Table) GetPlayers() ([]string, error) {

	t.mu.RLock()
	defer t.mu.RUnlock()

	seats := t.sm.GetSeats()

	players := make([]string, 0)
	for _, s := range seats {
		if s.Player == nil {
			continue
		}

		players = append(players, s.Player.(string))
	}

	return players, nil
}

func (t *Table) PrintState() {
	fmt.Printf("Table (id=%s, status=%d)\n",
		t.id,
		t.status,
	)
}

func (t *Table) OnPlayerJoined(fn func(string, int)) {
	t.onPlayerJoined = fn
}

func (t *Table) OnPlayerLeft(fn func(string, int)) {
	t.onPlayerLeft = fn
}

func (t *Table) OnPlayerDrained(fn func(string, int)) {
	t.onPlayerDrained = fn
}
