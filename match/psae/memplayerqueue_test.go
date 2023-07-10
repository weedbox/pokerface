package psae

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MemoryPlayerQueue(t *testing.T) {

	q := NewMemoryPlayerQueue()

	q.Publish(&Player{
		ID:   "test_player1",
		Name: "Test Player 1",
	})

	ch, err := q.Subscribe()
	assert.Nil(t, err)

	p := <-ch
	assert.Equal(t, "test_player1", p.ID)

	err = q.Close()
	assert.Nil(t, err)
}
