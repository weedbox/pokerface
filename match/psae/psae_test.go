package psae

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Engine_Join_AllocateTable(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*1)),
	)
	defer p.Close()

	// Preparing players
	for i := 0; i < 18; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 2)

	count, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}

func Test_Engine_Join_Redispatch(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*1)),
	)
	defer p.Close()

	// Preparing players
	for i := 0; i < 3; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 2)

	// It should not be enough to allocate a new table
	count, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	// Add more players
	for i := 0; i < 1; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 2)

	// players should be re-dispatched. with new players, one table was created
	count, err = p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func Test_Engine_AssertTableState(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	// Preparing table state
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

	err := p.AssertTableState(ts)
	assert.Nil(t, err)

	result, err := p.SeatMap().GetTableState(ts.ID)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.ID, ts.ID)
	assert.Equal(t, len(result.Players), len(ts.Players))
	assert.Equal(t, result.AvailableSeats, ts.AvailableSeats)
	assert.Equal(t, result.TotalSeats, ts.TotalSeats)
}

func Test_Engine_UpdateTableState(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	// Preparing table state
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

	err := p.AssertTableState(ts)
	assert.Nil(t, err)

	// New table state
	newts := &TableState{
		ID:             "test_table_1",
		Players:        make(map[string]*Player),
		AvailableSeats: 4,
		TotalSeats:     9,
	}

	for i := 0; i < 3; i++ {
		p := NewTestPlayer()
		newts.Players[p.ID] = p
	}

	// Update table
	newts, err = p.UpdateTableState(newts)
	assert.Nil(t, err)

	result, err := p.SeatMap().GetTableState(newts.ID)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Players), len(newts.Players))
	assert.Equal(t, result.AvailableSeats, newts.AvailableSeats)
	assert.Equal(t, result.TotalSeats, newts.TotalSeats)
}

func Test_Engine_BreakTable(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	// Preparing table state
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

	err := p.AssertTableState(ts)
	assert.Nil(t, err)

	err = p.BreakTable(ts.ID)
	assert.Nil(t, err)

	result, err := p.SeatMap().GetTableState(ts.ID)
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func Test_Engine_MatchPlayers(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	p1 := NewTestPlayer()
	p2 := NewTestPlayer()

	go func() {
		err := p.MatchPlayers([]*Player{
			p1,
			p2,
		})
		assert.Nil(t, err)
	}()

	ch, err := p.MatchQueue().Subscribe()
	assert.Nil(t, err)

	m := <-ch
	assert.Equal(t, p1.ID, m.Players[0].ID)
	assert.Equal(t, p2.ID, m.Players[1].ID)
}

func Test_Engine_DispatchPlayer(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	err := p.DispatchPlayer(&Player{
		ID:   "test_player_1",
		Name: "Test Player 1",
	})
	assert.Nil(t, err)

	ch, err := p.DispatchQueue().Subscribe()
	assert.Nil(t, err)

	player := <-ch
	assert.Equal(t, "test_player_1", player.ID)
}

func Test_Engine_ReleasePlayer(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	err := p.ReleasePlayer(&Player{
		ID:   "test_player_1",
		Name: "Test Player 1",
	})
	assert.Nil(t, err)

	ch, err := p.ReleaseQueue().Subscribe()
	assert.Nil(t, err)

	player := <-ch
	assert.Equal(t, "test_player_1", player.ID)
}

func Test_Engine_EnterWaitingRoom(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	player := &Player{
		ID:   "test_player_1",
		Name: "Test Player 1",
	}

	err := p.EnterWaitingRoom(player)
	assert.Nil(t, err)
}

func Test_Engine_LeaveWaitingRoom(t *testing.T) {

	rto := NewTestRuntimeOptions()

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	player := &Player{
		ID:   "test_player_1",
		Name: "Test Player 1",
	}

	err := p.EnterWaitingRoom(player)
	assert.Nil(t, err)

	err = p.LeaveWaitingRoom(player.ID)
	assert.Nil(t, err)
}

func Test_Engine_DrainWaitingRoom(t *testing.T) {

	drainedCount := 0
	done := make(chan struct{})

	rto := NewTestRuntimeOptions()
	rto.WaitingRoomDrained = func(p PSAE, player *Player) {
		drainedCount++
		assert.LessOrEqual(t, drainedCount, 7)

		if drainedCount == 7 {
			done <- struct{}{}
		}
	}

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	// Prepare players
	for i := 0; i < 7; i++ {
		player := NewTestPlayer()
		err := p.EnterWaitingRoom(player)
		assert.Nil(t, err)
	}

	err := p.DrainWaitingRoom()
	assert.Nil(t, err)

	<-done

	assert.Equal(t, 7, drainedCount)
}

func Test_Engine_FlushWaitingRoom(t *testing.T) {

	drainedCount := 0
	matchedCount := 0
	expected := map[int]int{
		1: 9,
		2: 9,
	}

	done := make(chan struct{})

	rto := NewTestRuntimeOptions()
	rto.WaitingRoomDrained = func(p PSAE, player *Player) {
		drainedCount++
		assert.LessOrEqual(t, drainedCount, 3)
		if drainedCount == 3 {
			done <- struct{}{}
		}
	}
	rto.WaitingRoomMatched = func(p PSAE, players []*Player) {
		matchedCount++
		assert.Equal(t, expected[matchedCount], len(players))
	}

	p := NewPSAE(
		WithRuntime(NewTestRuntime(rto)),
	)
	defer p.Close()

	// Prepare players
	for i := 0; i < 21; i++ {
		player := NewTestPlayer()
		err := p.EnterWaitingRoom(player)
		assert.Nil(t, err)
	}

	err := p.FlushWaitingRoom()
	assert.Nil(t, err)

	<-done

	assert.Equal(t, 3, drainedCount)
}
