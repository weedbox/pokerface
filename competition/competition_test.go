package competition

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface/table"
)

func Test_Competition_Basic(t *testing.T) {

	opts := NewOptions()
	opts.TableAllocationPeriod = 1

	tb := NewNativeTableBackend(table.NewNativeBackend())
	c := NewCompetition(
		opts,
		WithTableBackend(tb),
	)
	defer c.Close()

	assert.Nil(t, c.Start())

	// Registering
	assert.Nil(t, c.Register("player_1", 10000))
	assert.Nil(t, c.Register("player_2", 10000))
	assert.Nil(t, c.Register("player_3", 10000))
	assert.Nil(t, c.Register("player_4", 10000))
	assert.Nil(t, c.Register("player_5", 10000))
	assert.Nil(t, c.Register("player_6", 10000))
	assert.Nil(t, c.Register("player_7", 10000))

	players := c.GetPlayers()
	assert.Equal(t, 7, len(players))
	for _, p := range players {
		assert.True(t, p.Participated)
	}

	time.Sleep(2 * time.Second)

	// Allocated one table
	assert.Equal(t, 1, c.GetTableCount())
}
