package psae

import (
	"sync"
	"time"
)

type MemoryWaitingRoom struct {
	players *Stack
	timer   *time.Timer
	mu      sync.RWMutex
	wt      time.Duration
}

func NewMemoryWaitingRoom(waitingTime time.Duration) *MemoryWaitingRoom {

	r := &MemoryWaitingRoom{
		players: NewStack(),
		wt:      waitingTime,
	}

	r.timer = time.NewTimer(r.wt)
	r.timer.Stop()

	return r
}

func (mwr *MemoryWaitingRoom) Enter(p PSAE, player *Player) error {

	mwr.mu.Lock()
	defer mwr.mu.Unlock()

	mwr.players.Push(player)

	// Only one player means this room is empty beofre
	if mwr.players.Len() == 1 {

		// Setup timer to check this room
		mwr.timer.Reset(mwr.wt)

		go func() {
			select {
			case <-mwr.timer.C:
				mwr.Flush(p)
			}
		}()
	}

	p.EmitWaitingRoomEntered(player)

	return nil
}

func (mwr *MemoryWaitingRoom) Leave(p PSAE, pid string) error {

	mwr.mu.Lock()
	defer mwr.mu.Unlock()

	for e := mwr.players.List().Front(); e != nil; e = e.Next() {
		et := e.Value.(*Player)
		if et.ID == pid {
			mwr.players.List().Remove(e)
			break
		}
	}

	return nil
}

func (mwr *MemoryWaitingRoom) Drain(p PSAE) error {

	mwr.mu.Lock()
	defer mwr.mu.Unlock()

	players := mwr.players
	mwr.players = NewStack()

	for e := players.List().Front(); e != nil; e = e.Next() {
		player := e.Value.(*Player)
		p.EmitWaitingRoomDrained(player)
	}

	return nil
}

func (mwr *MemoryWaitingRoom) Match(p PSAE) error {

	mwr.mu.Lock()
	defer mwr.mu.Unlock()

	// More than minimum players
	if mwr.players.Len() < p.Game().MinInitialPlayers {
		return nil
	}

	if p.Game().TableLimit > -1 {

		tc, err := p.SeatMap().GetTableCount()
		if err != nil {
			return err
		}

		// Check if number of table have reached the maximum
		if tc >= p.Game().TableLimit {

			// Do nothing
			return nil
		}
	}

	players := make([]*Player, 0, p.Game().MaxPlayersPerTable)
	for v := mwr.players.Pop(); v != nil; v = mwr.players.Pop() {
		player := v.(*Player)
		players = append(players, player)

		// satisfy condition for new table
		if len(players) == p.Game().MaxPlayersPerTable {
			p.EmitWaitingRoomMatched(players)
			return nil
		}
	}

	p.EmitWaitingRoomMatched(players)

	return nil
}

func (mwr *MemoryWaitingRoom) Flush(p PSAE) error {

	if mwr.players.Len() == 0 {
		return nil
	}

	for mwr.players.Len() >= p.Game().MinInitialPlayers {
		err := mwr.Match(p)
		if err != nil {
			return err
		}
	}

	// Drain the rest of players in the waiting room
	if mwr.players.Len() > 0 {
		return mwr.Drain(p)
	}

	return nil
}
