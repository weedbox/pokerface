package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaitGroup(t *testing.T) {

	tm := NewTaskManager()

	wr := NewWaitReady("ready")
	wr.PreparePlayerStates(2)

	tm.AddTask(wr)

	assert.Equal(t, 1, tm.TaskCount())

	// No one is ready
	tm.Execute()
	assert.False(t, tm.IsCompleted())

	// Only one player was ready
	wr.Ready(0)
	tm.Execute()
	assert.False(t, tm.IsCompleted())

	// Second player was ready
	wr.Ready(1)
	tm.Execute()
	assert.True(t, wr.IsCompleted())
	assert.True(t, tm.IsCompleted())
}
