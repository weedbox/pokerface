package psae

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Stress_90_Players(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*2)),
	)
	defer p.Close()

	// Preparing 90 players
	for i := 0; i < 90; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 5)

	// Check table count
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 10, tableCount)
}

func Test_Stress_900_Players(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*2)),
	)
	defer p.Close()

	// Preparing 900 players
	for i := 0; i < 900; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 5)

	// Check table count
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 100, tableCount)
}

func Test_Stress_9000_Players(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*2)),
	)
	defer p.Close()

	// Preparing 9000 players
	for i := 0; i < 9000; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 5)

	// Check table count
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 1000, tableCount)
}

func Test_Stress_90000_Players(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*2)),
	)
	defer p.Close()

	// Preparing 90000 players
	for i := 0; i < 90000; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 5)

	// Check table count
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 10000, tableCount)
}

func Test_Stress_RampUp(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*2)),
	)
	defer p.Close()

	// Preparing 10 tables
	for i := 0; i < 90; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 3)

	// Check table count
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 10, tableCount)

	// Preparing 10 tables
	for i := 0; i < 90; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 3)

	// Check table count
	tableCount, err = p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 20, tableCount)

}

func Test_Stress_RampUp_PlayerLeft(t *testing.T) {

	b := NewBackend()
	b.AllocateTable = func() (*TableState, error) {
		return NewTestTableState(0), nil
	}

	p := NewPSAE(
		WithBackend(b),
		WithWaitingRoom(NewMemoryWaitingRoom(time.Second*2)),
	)
	defer p.Close()

	// Preparing 10 tables
	for i := 0; i < 90; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 3)

	// Check table count
	tableCount, err := p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 10, tableCount)

	tables, err := p.SeatMap().GetAllTables()
	assert.Nil(t, err)

	// Remove one player from each table
	for _, t := range tables {

		for pid, _ := range t.Players {
			delete(t.Players, pid)
			t.AvailableSeats--
			break
		}
	}

	// Preparing 10 tables
	for i := 0; i < 90; i++ {
		player := NewTestPlayer()
		err := p.Join(player)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second * 3)

	// Check table count
	tableCount, err = p.SeatMap().GetTableCount()
	assert.Nil(t, err)
	assert.Equal(t, 20, tableCount)

}
