package match

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestTableState struct {
	PlayerCount int
	Seats       map[int]string
}

func Test_Match_Basic(t *testing.T) {

	opts := NewOptions(uuid.New().String())
	opts.WaitingPeriod = 1

	m := NewMatch(opts)
	defer m.Close()

	totalPlayers := 900
	for i := 0; i < totalPlayers; i++ {
		assert.Nil(t, m.Register(fmt.Sprintf("player_%d", i+1)))
	}

	time.Sleep(2 * time.Second)

	// Check tables
	assert.Equal(t, int64(totalPlayers/opts.MaxSeats), m.TableMap().Count())

	tables, err := m.TableMap().GetTables()
	assert.Nil(t, err)
	assert.Equal(t, totalPlayers/opts.MaxSeats, len(tables))

	// Check players of table
	playerCount := 0
	for _, table := range tables {
		playerCount += table.GetPlayerCount()
		assert.Equal(t, 9, table.GetPlayerCount())
	}

	assert.Equal(t, totalPlayers, playerCount)

	// no player stay in the waiting room
	countOfRoom, _ := m.WaitingRoom().Count()
	assert.Zero(t, countOfRoom)
}

func Test_Match_BreakTable_UntilOnlyOneTable(t *testing.T) {

	totalPlayers := 900
	expectedPlayers := 6
	removed := 0
	brokenCount := 0
	tb := NewDummyTableBackend()

	removePlayers := func(m Match, removeCount int) bool {

		var table *Table

		// Find a table with enough players to be removed
		tables := tb.getTables()
		for tableID, ts := range tables {

			if ts.GetStatus() == TableStatus_Busy {
				continue
			}

			playerCount := ts.GetPlayerCount()

			if playerCount < 2 || playerCount <= removeCount {
				//t.Logf("table_id=%s, playerCount=%d", tableID, playerCount)
				continue
			}

			found, err := tb.GetTable(tableID)
			if err != nil {
				continue
			}

			table = found
			break
		}

		if table == nil {
			return false
		}

		// Getting local table states
		seats := table.SeatManager().GetSeats()

		// Preparing seat changes
		sc := NewSeatChanges()
		count := 0
		for _, s := range seats {

			if s.Player == nil {
				continue
			}

			// Remove players
			sc.Seats[s.ID] = "left"

			count++
			if count == removeCount {
				break
			}
		}

		assert.Equal(t, removeCount, count)
		removed += removeCount

		// Apply changes on local table state
		assert.Nil(t, table.ApplySeatChanges(sc))

		t.Logf("[%d] (%d/%d) Removed %d players (table_id=%s, left=%d, status=%d, pending=%d)",
			removed,
			totalPlayers-removed,
			totalPlayers,
			removeCount,
			table.ID(),
			table.GetPlayerCount(),
			table.GetStatus(),
			m.Dispatcher().GetPendingCount(),
		)

		// Emit event to simulate table state changes
		tb.UpdateTable(table.ID(), sc)

		return true
	}

	opts := NewOptions(uuid.New().String())
	opts.WaitingPeriod = 1

	m := NewMatch(
		opts,
		WithTableBackend(tb),
		WithPlayerJoinedCallback(func(m Match, table *Table, seatID int, playerID string) {
			table, err := tb.GetTable(table.ID())
			assert.Nil(t, err)

			// Confirm that the players are in their respective seats
			seat := table.SeatManager().GetSeat(seatID)
			assert.Equal(t, playerID, seat.Player.(string))
		}),
		WithTableBrokenCallback(func(m Match, table *Table) {

			brokenCount++
			t.Logf("[Break %d] Break table (table_id=%s, left=%d, status=%d)",
				brokenCount,
				table.ID(),
				table.GetPlayerCount(),
				table.GetStatus(),
			)
		}),
	)
	defer m.Close()

	// Preparing players and tables
	for i := 0; i < totalPlayers; i++ {
		assert.Nil(t, m.Register(fmt.Sprintf("player_%d", i+1)))
	}

	time.Sleep(time.Second)

	totalTables := m.TableMap().Count()

	// Randomly select a table and remove a player each time
	for i := 0; i < totalPlayers-expectedPlayers; i++ {
		for !removePlayers(m, 1) {
			// Attempt to remove player in one second
			time.Sleep(time.Second)
		}
	}

	time.Sleep(1 * time.Second)

	// It should be two tables left
	assert.Equal(t, totalTables-int64(brokenCount), m.TableMap().Count())
	if !assert.Equal(t, int64(1), m.TableMap().Count()) {
		m.PrintTables()
	}

	// Check total players
	assert.Equal(t, int64(totalPlayers-removed), m.GetPlayerCount())
	assert.Equal(t, int64(expectedPlayers), m.GetPlayerCount())
}
