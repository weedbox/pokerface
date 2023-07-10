package psae

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MemorySeatMap_CreateTable(t *testing.T) {

	msm := NewMemorySeatMap()

	ts := &TableState{
		ID:             "test_table",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 10; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)
	assert.Len(t, msm.tables, 1)
	assert.Equal(t, 1, msm.ordered.Len())
}

func Test_MemorySeatMap_CreateTable_Ordering(t *testing.T) {

	msm := NewMemorySeatMap()

	// Table 1
	ts := &TableState{
		ID:             "test_table_1",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 2; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)
	assert.Equal(t, 1, msm.ordered.Len())

	// Table 2
	ts2 := &TableState{
		ID:             "test_table_2",
		Players:        make(map[string]*Player),
		AvailableSeats: 1,
		TotalSeats:     9,
	}

	for i := 0; i < 9; i++ {
		p := NewTestPlayer()
		ts2.Players[p.ID] = p
	}

	err = msm.CreateTable(ts2)
	assert.Nil(t, err)
	assert.Equal(t, 2, msm.ordered.Len())

	// Table 3
	ts3 := &TableState{
		ID:             "test_table_3",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 4; i++ {
		p := NewTestPlayer()
		ts3.Players[p.ID] = p
	}

	err = msm.CreateTable(ts3)
	assert.Nil(t, err)
	assert.Equal(t, 3, msm.ordered.Len())

	// Check for ordering
	lastTable := msm.ordered.Front().Value.(*TableState)
	for e := msm.ordered.Front(); e != nil; e = e.Next() {
		et := e.Value.(*TableState)

		assert.GreaterOrEqual(t, len(et.Players), len(lastTable.Players))
		lastTable = et
	}
}

func Test_MemorySeatMap_DestroyTable(t *testing.T) {

	msm := NewMemorySeatMap()

	ts := &TableState{
		ID:             "test_table",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 10; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)
	assert.Len(t, msm.tables, 1)
	assert.Equal(t, 1, msm.ordered.Len())

	err = msm.DestroyTable(ts.ID)
	assert.Nil(t, err)
	assert.Len(t, msm.tables, 0)
	assert.Equal(t, 0, msm.ordered.Len())
}

func Test_MemorySeatMap_GetTable(t *testing.T) {

	msm := NewMemorySeatMap()

	ts := &TableState{
		ID:             "test_table",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 10; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)

	received, err := msm.GetTableState(ts.ID)
	assert.Nil(t, err)
	assert.Equal(t, received.ID, ts.ID)
	assert.Len(t, received.Players, len(ts.Players))
	assert.Equal(t, received.AvailableSeats, ts.AvailableSeats)
}

func Test_MemorySeatMap_UpdateTable(t *testing.T) {

	msm := NewMemorySeatMap()

	ts := &TableState{
		ID:             "test_table",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 2; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)

	newts := &TableState{
		ID:             "test_table",
		Players:        make(map[string]*Player),
		AvailableSeats: 4,
		TotalSeats:     9,
	}

	for i := 0; i < 3; i++ {
		p := NewTestPlayer()
		newts.Players[p.ID] = p
	}

	changed, err := msm.UpdateTable(newts)
	assert.Nil(t, err)

	received, err := msm.GetTableState(ts.ID)
	assert.Nil(t, err)
	assert.Equal(t, received.ID, changed.ID)
	assert.Len(t, received.Players, len(newts.Players))
	assert.Equal(t, received.AvailableSeats, newts.AvailableSeats)
}

func Test_MemorySeatMap_UpdateTable_Ordering(t *testing.T) {

	msm := NewMemorySeatMap()

	// Table 1
	ts := &TableState{
		ID:             "test_table_1",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 2; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)

	// Table 2
	ts2 := &TableState{
		ID:             "test_table_2",
		Players:        make(map[string]*Player),
		AvailableSeats: 1,
		TotalSeats:     9,
	}

	for i := 0; i < 9; i++ {
		p := NewTestPlayer()
		ts2.Players[p.ID] = p
	}

	err = msm.CreateTable(ts2)
	assert.Nil(t, err)

	// Table 3
	ts3 := &TableState{
		ID:             "test_table_3",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 4; i++ {
		p := NewTestPlayer()
		ts3.Players[p.ID] = p
	}

	err = msm.CreateTable(ts3)
	assert.Nil(t, err)

	// Add 3 players to test_table_1
	newts := &TableState{
		ID:             "test_table_1",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 5; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	_, err = msm.UpdateTable(newts)
	assert.Nil(t, err)

	// Check for ordering
	lastTable := msm.ordered.Front().Value.(*TableState)
	for e := msm.ordered.Front(); e != nil; e = e.Next() {
		et := e.Value.(*TableState)
		assert.GreaterOrEqual(t, len(et.Players), len(lastTable.Players))
		lastTable = et
	}
}

func Test_MemorySeatMap_FindAvailableTable(t *testing.T) {

	msm := NewMemorySeatMap()

	// Table 1
	ts := &TableState{
		ID:             "test_table_1",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 2; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)

	// Table 2
	ts2 := &TableState{
		ID:             "test_table_2",
		Players:        make(map[string]*Player),
		AvailableSeats: 1,
		TotalSeats:     9,
	}

	for i := 0; i < 8; i++ {
		p := NewTestPlayer()
		ts2.Players[p.ID] = p
	}

	err = msm.CreateTable(ts2)
	assert.Nil(t, err)

	// Table 3
	ts3 := &TableState{
		ID:             "test_table_3",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 4; i++ {
		p := NewTestPlayer()
		ts3.Players[p.ID] = p
	}

	err = msm.CreateTable(ts3)
	assert.Nil(t, err)

	// Highest number of players
	availableTable, err := msm.FindAvailableTable(&TableCondition{
		HighestNumberOfPlayers: true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, availableTable)
	assert.Len(t, availableTable.Players, 8)

	// Lowest number of players
	availableTable, err = msm.FindAvailableTable(&TableCondition{
		HighestNumberOfPlayers: false,
	})
	assert.Nil(t, err)
	assert.NotNil(t, availableTable)
	assert.Len(t, availableTable.Players, 2)
}

func Test_MemorySeatMap_TotalPlayers(t *testing.T) {

	msm := NewMemorySeatMap()

	ts := &TableState{
		ID:             "test_table",
		Players:        make(map[string]*Player),
		AvailableSeats: 5,
		TotalSeats:     9,
	}

	for i := 0; i < 9; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	err := msm.CreateTable(ts)
	assert.Nil(t, err)

	count, err := msm.GetTotalPlayers()
	assert.Nil(t, err)
	assert.Equal(t, int64(9), count)

	// Replace all players
	players := make(map[string]*Player)
	for i := 0; i < 5; i++ {
		p := NewTestPlayer()
		players[p.ID] = p
	}

	ts.Players = players

	_, err = msm.UpdateTable(ts)
	assert.Nil(t, err)

	count, err = msm.GetTotalPlayers()
	assert.Nil(t, err)
	assert.Equal(t, int64(5), count)

	// Destroy table
	err = msm.DestroyTable(ts.ID)
	assert.Nil(t, err)

	count, err = msm.GetTotalPlayers()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
}
