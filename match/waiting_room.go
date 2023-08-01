package match

import (
	"sync"
	"time"

	"github.com/weedbox/timebank"
)

type WaitingRoom interface {
	Count() (int, error)
	Enter(playerID string) error
	Leave(playerID string) error
	Drain() error
	Match() error
	Flush() error
}

type waitingRoom struct {
	m       Match
	players *Stack
	tb      *timebank.TimeBank
	mu      sync.RWMutex
}

func NewWaitingRoom(m Match) WaitingRoom {

	wr := &waitingRoom{
		m:       m,
		players: NewStack(),
		tb:      timebank.NewTimeBank(),
	}

	return wr
}

func (wr *waitingRoom) Count() (int, error) {
	return wr.players.Len(), nil
}

func (wr *waitingRoom) Enter(playerID string) error {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	// First player is coming
	if wr.players.Len() == 0 {

		// Setup timer to check this room later
		wr.tb.NewTask(time.Duration(wr.m.Options().WaitingPeriod)*time.Second, func(isCancelled bool) {

			if isCancelled {
				return
			}

			wr.Flush()
		})
	}

	wr.players.Push(playerID)

	// Dispatch players immediately
	if wr.players.Len() >= wr.m.Options().MaxSeats {
		return wr.flush()
	}

	return nil
}

func (wr *waitingRoom) Leave(playerID string) error {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	for e := wr.players.List().Front(); e != nil; e = e.Next() {
		pid := e.Value.(string)
		if pid == playerID {
			wr.players.List().Remove(e)
			break
		}
	}

	return nil
}

func (wr *waitingRoom) Drain() error {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	return wr.drain()
}

func (wr *waitingRoom) Match() error {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	return wr.match()
}

func (wr *waitingRoom) Flush() error {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	return wr.flush()
}

func (wr *waitingRoom) drain() error {

	players := wr.players

	// Clean waiting room
	wr.players = NewStack()

	// drain all players from the waiting room
	ids := make([]string, 0)
	for e := players.List().Front(); e != nil; e = e.Next() {
		playerID := e.Value.(string)
		ids = append(ids, playerID)
	}

	wr.m.Runner().DrainWaitingRoomPlayers(wr.m, ids)

	return nil
}

func (wr *waitingRoom) match() error {

	// Less than minimum initial players
	if wr.players.Len() < wr.m.Options().MinInitialPlayers {
		return nil
	}

	if wr.m.Options().MaxTables > -1 {

		// Check if the number of table have reached the maximum
		if wr.m.TableMap().Count() >= int64(wr.m.Options().MaxTables) {

			// Do nothing
			return nil
		}
	}

	// Dispatch players
	players := make([]string, 0, wr.m.Options().MaxSeats)
	for v := wr.players.Pop(); v != nil; v = wr.players.Pop() {
		playerID := v.(string)
		players = append(players, playerID)

		// satisfy condition for new table
		if len(players) == wr.m.Options().MaxSeats {
			break
		}
	}

	// Allocate table for players
	wr.m.AllocateTableWithPlayers(players)

	return nil
}

func (wr *waitingRoom) flush() error {

	wr.tb.Cancel()

	for wr.players.Len() >= wr.m.Options().MinInitialPlayers {
		err := wr.match()
		if err != nil {
			return err
		}
	}

	// Drain the rest of players in the waiting room
	if wr.players.Len() > 0 {
		return wr.drain()
	}

	return nil
}
