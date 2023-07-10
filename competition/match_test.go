package competition

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface/table"
)

func Test_Match_AllocateNewTable(t *testing.T) {

	opts := NewOptions()
	opts.MaxTables = -1
	opts.TableAllocationPeriod = 1
	opts.Table.InitialPlayers = 5

	tb := NewNativeTableBackend(table.NewNativeBackend())
	c := NewCompetition(
		opts,
		WithTableBackend(tb),
	)

	assert.Nil(t, c.Start())
	assert.Equal(t, 0, c.GetTableCount())

	// Registering
	c.Register("player_1", 10000)
	c.Register("player_2", 10000)
	c.Register("player_3", 10000)
	c.Register("player_4", 10000)
	c.Register("player_5", 10000)

	players := c.GetPlayers()
	assert.Equal(t, 5, len(players))
	for _, p := range players {
		assert.True(t, p.Participated)
	}

	time.Sleep(time.Second * 2)

	// It should allocate a new table for players
	assert.Equal(t, 1, c.GetTableCount())
}

func Test_Match_AllocateMoreTables(t *testing.T) {

	opts := NewOptions()
	opts.MaxTables = -1
	opts.TableAllocationPeriod = 1
	opts.Table.InitialPlayers = 5
	opts.Table.MaxSeats = 9

	tb := NewNativeTableBackend(table.NewNativeBackend())
	c := NewCompetition(
		opts,
		WithTableBackend(tb),
	)

	assert.Nil(t, c.Start())
	assert.Equal(t, 0, c.GetTableCount())

	// Registering
	totalPlayers := 100
	for i := 0; i < totalPlayers; i++ {
		c.Register(fmt.Sprintf("player_%d", i+1), 10000)
	}

	players := c.GetPlayers()
	assert.Equal(t, totalPlayers, len(players))
	for _, p := range players {
		assert.True(t, p.Participated)
	}

	time.Sleep(time.Second * 2)

	// It should allocate more tables for players
	assert.Equal(t, totalPlayers/opts.Table.MaxSeats, c.GetTableCount())
}
