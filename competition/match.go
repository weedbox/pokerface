package competition

import (
	"fmt"

	"github.com/weedbox/pokerface/match"
)

type NativeMatchTableBackend struct {
	c              Competition
	onTableUpdated func(tableID string, sc *match.SeatChanges)
}

func NewNativeMatchTableBackend(c Competition) *NativeMatchTableBackend {

	nmtb := &NativeMatchTableBackend{
		c:              c,
		onTableUpdated: func(tableID string, sc *match.SeatChanges) {},
	}

	return nmtb
}

func (nmtb *NativeMatchTableBackend) Allocate(maxSeats int) (*match.Table, error) {

	// Create a new table
	ts, err := nmtb.c.TableManager().CreateTable()
	if err != nil {
		return nil, err
	}

	// Preparing table state
	t := match.NewTable(maxSeats)
	t.SetID(ts.ID)

	fmt.Printf("Allocated Table (id=%s, seats=%d)\n", ts.ID, maxSeats)

	return t, nil
}

func (nmtb *NativeMatchTableBackend) Release(tableID string) error {
	return nmtb.c.TableManager().ReleaseTable(tableID)
}

func (nmtb *NativeMatchTableBackend) Activate(tableID string) error {

	err := nmtb.c.TableManager().ActivateTable(tableID)
	if err != nil {
		return err
	}

	fmt.Printf("Activated Table (id=%s)\n", tableID)

	return nil
}

func (nmtb *NativeMatchTableBackend) Reserve(tableID string, seatID int, playerID string) error {

	//fmt.Printf("<= Reseve Seat (table_id=%s, seat=%d, player=%s)\n", tableID, seatID, playerID)

	_, err := nmtb.c.ReserveSeat(tableID, seatID, playerID)
	if err != nil {
		return err
	}

	return nil
}

func (nmtb *NativeMatchTableBackend) GetTable(tableID string) (*match.Table, error) {

	ts := nmtb.c.TableManager().GetTableState(tableID)

	t := match.NewTable(ts.Options.MaxSeats)
	t.SetID(ts.ID)

	return t, nil
}

func (nmtb *NativeMatchTableBackend) UpdateTable(tableID string, sc *match.SeatChanges) error {
	nmtb.onTableUpdated(tableID, sc)
	return nil
}

func (nmtb *NativeMatchTableBackend) OnTableUpdated(fn func(tableID string, sc *match.SeatChanges)) {
	nmtb.onTableUpdated = fn
}
