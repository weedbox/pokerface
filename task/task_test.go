package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {

	tm := NewTaskManager()

	task1 := NewTask(func(ct *CustomizedTask) bool {
		return true
	})

	tm.AddTask(task1)

	task2 := NewTask(func(ct *CustomizedTask) bool {
		return true
	})

	tm.AddTask(task2)

	assert.Equal(t, 2, tm.TaskCount())
	assert.False(t, tm.IsCompleted())

	tm.Execute()

	assert.True(t, tm.IsCompleted())
}
