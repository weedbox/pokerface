package psae

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Balance_PlayerCount_LessThanOrEqualToThree(t *testing.T) {

	p := NewPSAE()
	defer p.Close()

	// Preparing a full table
	ts := NewTestTableState(0)
	ts.AvailableSeats = 0

	newPlayers := make(map[string]*Player)
	for i := 0; i < 9; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p

		if i < 3 {
			newPlayers[p.ID] = p
		}
	}

	err := p.AssertTableState(ts)
	assert.Nil(t, err)

	// Preparing a table which is not full
	nfts := NewTestTableState(5)
	err = p.AssertTableState(nfts)
	assert.Nil(t, err)

	// Update table to simulate player leaves
	ts.Players = newPlayers
	ts.AvailableSeats = 6
	ts.Status = TableStatus_Suspend
	updated, err := p.UpdateTableState(ts)
	assert.Nil(t, err)
	assert.Equal(t, TableStatus_Broken, updated.Status)

	time.Sleep(time.Second * 2)

	// first table should be destroyed
	oldts, err := p.SeatMap().GetTableState(ts.ID)
	assert.Nil(t, err)
	assert.Nil(t, oldts)

	// second table should be updated
	newts, err := p.SeatMap().GetTableState(nfts.ID)
	assert.Nil(t, err)
	assert.NotNil(t, newts)

	// The table will be filled
	assert.Equal(t, 1, newts.AvailableSeats)
	assert.Equal(t, 8, len(newts.Players))
}

func Test_Balance_PlayerCount_LessThanAverage(t *testing.T) {

	p := NewPSAE()
	defer p.Close()

	// Create three tables which contains 9 players
	for i := 0; i < 3; i++ {
		ts := NewTestTableState(9)
		err := p.AssertTableState(ts)
		assert.Nil(t, err)
	}

	// Create three tables which contains 5 players
	for i := 0; i < 3; i++ {
		ts := NewTestTableState(5)
		err := p.AssertTableState(ts)
		assert.Nil(t, err)
	}

	// Add a table it should be broken
	ts := NewTestTableState(4)
	err := p.AssertTableState(ts)
	assert.Nil(t, err)

	// Should be 7 tables
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 7, tableCount)

	// Update table states 10 times
	for i := 0; i < 10; i++ {
		ts.LastGameID = fmt.Sprintf("%d", i)
		_, err := p.UpdateTableState(ts)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 2)

	// the last table should be destroyed
	tableCount, err = p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 6, tableCount)

	// Check the number of players
	tables, err := p.SeatMap().GetAllTables()
	assert.Nil(t, err)
	playerCount := 0
	for _, t := range tables {
		playerCount += len(t.Players)
	}

	totalPlayer, err := p.SeatMap().GetTotalPlayers()
	assert.Nil(t, err)
	assert.Equal(t, int64(playerCount), totalPlayer)
}

func Test_Balance_PlayerCount_OnePlayerLeft(t *testing.T) {

	p := NewPSAE()
	defer p.Close()

	// Create three tables which contains 9 players
	for i := 0; i < 2; i++ {
		ts := NewTestTableState(9)
		err := p.AssertTableState(ts)
		assert.Nil(t, err)
	}

	// Create three tables which contains 5 players
	for i := 0; i < 2; i++ {
		ts := NewTestTableState(5)
		err := p.AssertTableState(ts)
		assert.Nil(t, err)
	}

	// Add a table it should be broken
	ts := NewTestTableState(1)
	err := p.AssertTableState(ts)
	assert.Nil(t, err)

	// Should be 5 tables
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 5, tableCount)

	// Update table states to trigger balancing
	ts.LastGameID = fmt.Sprintf("%d", 0)
	ts, err = p.UpdateTableState(ts)
	assert.Nil(t, err)
	assert.Equal(t, ts.Status, TableStatus_Broken)

	time.Sleep(time.Second * 2)

	// the last table should be destroyed
	tableCount, err = p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 4, tableCount)

	// Check the number of players
	tables, err := p.SeatMap().GetAllTables()
	assert.Nil(t, err)
	playerCount := 0
	for _, t := range tables {
		playerCount += len(t.Players)
	}

	totalPlayer, err := p.SeatMap().GetTotalPlayers()
	assert.Nil(t, err)
	assert.Equal(t, int64(playerCount), totalPlayer)
}

func Test_Balance_PlayerCount_TwoTablesShouldBeBroken(t *testing.T) {

	p := NewPSAE()
	defer p.Close()

	// Create three tables which contains 9 players
	for i := 0; i < 2; i++ {
		ts := NewTestTableState(9)
		err := p.AssertTableState(ts)
		assert.Nil(t, err)
	}

	// Create three tables which contains 5 players
	for i := 0; i < 2; i++ {
		ts := NewTestTableState(5)
		err := p.AssertTableState(ts)
		assert.Nil(t, err)
	}

	// Add two tables that should be broken
	tables := make([]*TableState, 0)
	for i := 0; i < 2; i++ {
		ts := NewTestTableState(2)
		err := p.AssertTableState(ts)
		assert.Nil(t, err)

		tables = append(tables, ts)
	}

	// Should be 5 tables
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 6, tableCount)

	for i := 0; i < 1; i++ {
		for _, ts := range tables {
			// Update table states to trigger balancing
			ts.LastGameID = fmt.Sprintf("%d", i)
			ts, err = p.UpdateTableState(ts)
			assert.Nil(t, err)
			assert.Equal(t, ts.Status, TableStatus_Broken)
		}
	}

	time.Sleep(time.Second * 2)

	// the last table should be destroyed
	tableCount, err = p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 4, tableCount)

	// Check the number of players
	tables, err = p.SeatMap().GetAllTables()
	assert.Nil(t, err)
	playerCount := 0
	for _, ts := range tables {
		t.Log(len(ts.Players))
		playerCount += len(ts.Players)
	}

	totalPlayer, err := p.SeatMap().GetTotalPlayers()
	assert.Nil(t, err)
	assert.Equal(t, int64(playerCount), totalPlayer)
}
