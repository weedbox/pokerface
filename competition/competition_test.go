package competition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface/table"
)

func Test_Competition_Basic(t *testing.T) {

	opts := NewOptions()

	tb := NewNativeTableBackend(table.NewNativeBackend())
	c := NewCompetition(
		opts,
		WithTableBackend(tb),
	)

	assert.Nil(t, c.Start())
	assert.Equal(t, 1, c.GetTableCount())

	// Registering
	c.Register("player_1", 10000)
	c.Register("player_2", 10000)
	c.Register("player_3", 10000)

	players := c.GetPlayers()
	assert.Equal(t, 3, len(players))
	for _, p := range players {
		assert.True(t, p.Participated)
	}
}
