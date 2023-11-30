package match

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_WaitingRoom_Basic(t *testing.T) {

	opts := NewOptions(uuid.New().String())
	opts.WaitingPeriod = 1

	m := NewMatch(opts)
	defer m.Close()

	totalPlayers := 900
	for i := 0; i < totalPlayers; i++ {
		assert.Nil(t, m.WaitingRoom().Enter(fmt.Sprintf("player_%d", i+1)))
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
		playerCount += table.SeatManager().GetPlayerCount()
		assert.Equal(t, 9, table.SeatManager().GetPlayerCount())
	}

	assert.Equal(t, totalPlayers, playerCount)

	// no player stay in the waiting room
	countOfRoom, _ := m.WaitingRoom().Count()
	assert.Zero(t, countOfRoom)
}

func Test_WaitingRoom_WaitingPeriod(t *testing.T) {

	opts := NewOptions(uuid.New().String())
	opts.WaitingPeriod = 5

	m := NewMatch(opts)
	defer m.Close()

	totalPlayers := 95
	for i := 0; i < totalPlayers; i++ {
		assert.Nil(t, m.WaitingRoom().Enter(fmt.Sprintf("player_%d", i+1)))
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
		playerCount += table.SeatManager().GetPlayerCount()
	}

	assert.Equal(t, totalPlayers-5, playerCount)

	// it should be 5 players stay in the waiting room
	countOfRoom, _ := m.WaitingRoom().Count()
	assert.Equal(t, 5, countOfRoom)

	// Waiting for final match
	time.Sleep(4 * time.Second)

	// Check tables
	assert.Equal(t, int64(1+(totalPlayers/opts.MaxSeats)), m.TableMap().Count())

	tables, err = m.TableMap().GetTables()
	assert.Nil(t, err)
	assert.Equal(t, 1+(totalPlayers/opts.MaxSeats), len(tables))

	// Check players of table
	playerCount = 0
	for _, table := range tables {
		playerCount += table.SeatManager().GetPlayerCount()
	}

	assert.Equal(t, totalPlayers, playerCount)

	// no player stay in the waiting room
	countOfRoom, _ = m.WaitingRoom().Count()
	assert.Zero(t, countOfRoom)
}

func Test_WaitingRoom_Pump_Up(t *testing.T) {

	opts := NewOptions(uuid.New().String())
	opts.MaxSeats = 9
	opts.MinInitialPlayers = 4
	opts.WaitingPeriod = 1

	m := NewMatch(opts)
	defer m.Close()

	totalPlayers := 900

	var wg sync.WaitGroup
	wg.Add(1)
	go func(total int) {

		unitPerLevel := 104

		for i := 0; i < total; {
			for j := 0; j < unitPerLevel; j++ {
				assert.Nil(t, m.WaitingRoom().Enter(fmt.Sprintf("player_%d", i+1)))

				i++
				if i == total {
					break
				}
			}

			t.Log("Entered", i)

			// Waiting for period of waiting room
			time.Sleep(2 * time.Second)
		}

		wg.Done()

	}(totalPlayers)

	wg.Wait()

	// Check tables
	tables, err := m.TableMap().GetTables()
	assert.Nil(t, err)

	// Check players of table
	playerCount := 0
	for _, table := range tables {
		playerCount += table.SeatManager().GetPlayerCount()
	}
	assert.Equal(t, totalPlayers, playerCount)

	// no player stay in the waiting room
	countOfRoom, _ := m.WaitingRoom().Count()
	assert.Zero(t, countOfRoom)
}
