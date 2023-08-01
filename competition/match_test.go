package competition

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface/table"
)

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
	defer c.Close()

	assert.Nil(t, c.Start())
	assert.Equal(t, int64(0), c.GetTableCount())

	// Registering
	totalPlayers := 101
	for i := 0; i < totalPlayers; i++ {
		t.Log(i + 1)
		assert.Nil(t, c.Register(fmt.Sprintf("player_%d", i+1), 10000))
	}

	players := c.GetPlayers()
	assert.Equal(t, totalPlayers, len(players))
	for _, p := range players {
		assert.True(t, p.Participated)
	}

	time.Sleep(time.Second * 2)

	// It should allocate more tables for players
	assert.Equal(t, int64(totalPlayers/opts.Table.MaxSeats), c.GetTableCount())

	// Two players are still waiting
	waitingCount, _ := c.Match().WaitingRoom().Count()
	assert.Equal(t, 2, waitingCount)
}
