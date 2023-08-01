package competition

import (
	"github.com/weedbox/pokerface/match"
	"github.com/weedbox/pokerface/table"
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

	c.TableBackend().OnTableUpdated(func(ts *table.State) {
		nmtb.EmitTableUpdated(ts)
	})

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

	// Activate immediately
	err = nmtb.c.TableManager().ActivateTable(ts.ID)
	if err != nil {
		return t, err
	}

	//fmt.Printf("Allocated Table (id=%s, seats=%d)\n", ts.ID, maxSeats)

	return t, nil
}

func (nmtb *NativeMatchTableBackend) Release(tableID string) error {
	return nmtb.c.TableManager().ReleaseTable(tableID)
}

func (nmtb *NativeMatchTableBackend) Reserve(tableID string, seatID int, playerID string) error {

	//fmt.Printf("Reseve Seat (table_id=%s, seat=%d, player=%s)\n", tableID, seatID, playerID)

	_, err := nmtb.c.ReserveSeat(tableID, -1, playerID)
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

func (nmtb *NativeMatchTableBackend) EmitTableUpdated(newts *table.State) {

	if newts.GameState == nil {
		return
	}

	//newts.PrintState()

	// Getting original table state
	ts := nmtb.c.TableManager().GetTableState(newts.ID)
	if ts == nil {
		return
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
	for _, p := range ts.Players {

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
	nmtb.onTableUpdated(ts.ID, sc)
}
func (nmtb *NativeMatchTableBackend) OnTableUpdated(fn func(tableID string, sc *match.SeatChanges)) {
	nmtb.onTableUpdated = fn
}
