package competition

import (
	"fmt"

	"github.com/weedbox/pokerface/match"
	"github.com/weedbox/pokerface/match/psae"
)

type MatchBackend interface {
	AllocateTable() (string, error)
	BreakTable(tableID string) error
	Join(tableID string, players []*psae.Player) ([]int, error)
	DispatchPlayer(playerID string) error
}

type NativeMatchBackend struct {
	c Competition
	m match.Match
}

func NewNativeMatchBackend(c Competition) MatchBackend {

	nmb := &NativeMatchBackend{
		c: c,
	}

	// Initializing match mechanism
	opts := match.NewOptions()
	opts.WaitingPeriod = c.GetOptions().TableAllocationPeriod
	opts.MaxTables = c.GetOptions().MaxTables
	opts.MaxSeats = c.GetOptions().Table.MaxSeats

	nmb.m = match.NewMatch(opts, match.WithBackend(nmb))

	return nmb
}

func (nmb *NativeMatchBackend) AllocateTable() (string, error) {

	// Create a new table
	ts, err := nmb.c.GetTableManager().CreateTable()
	if err != nil {
		return "", err
	}

	fmt.Printf("AllocateTable (id=%s)\n", ts.ID)

	// Activate immediately
	err = nmb.c.GetTableManager().ActivateTable(ts.ID)
	if err != nil {
		return ts.ID, err
	}

	return ts.ID, nil
}

func (nmb *NativeMatchBackend) BreakTable(tableID string) error {

	fmt.Printf("BreakTable (id=%s)\n", tableID)

	return nmb.c.GetTableManager().BreakTable(tableID)
}

func (nmb *NativeMatchBackend) Join(tableID string, players []*psae.Player) ([]int, error) {

	seats := make([]int, 0)

	// Attempt to reserve seats for players
	for _, p := range players {

		seatID, err := nmb.c.GetTableManager().ReserveSeat(tableID, -1, &PlayerInfo{
			ID: p.ID,
		})

		seats = append(seats, seatID)

		if err != nil {
			return seats, err
		}
	}

	return seats, nil
}

func (nmb *NativeMatchBackend) DispatchPlayer(playerID string) error {
	return nmb.m.Join(playerID)
}
