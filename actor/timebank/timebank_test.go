package timebank

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {

	tb := NewTimeBank()

	var wg sync.WaitGroup
	wg.Add(1)

	startTime := time.Now()

	err := tb.NewTask(3*time.Second, func(isCancelled bool) {
		wg.Done()
	})

	assert.Nil(t, err)

	wg.Wait()

	duration := time.Now().Unix() - startTime.Unix()
	assert.Equal(t, int64(3), duration)
}

func TestCancel(t *testing.T) {

	tb := NewTimeBank()

	var wg sync.WaitGroup
	wg.Add(1)

	err := tb.NewTask(3*time.Second, func(isCancelled bool) {
		assert.Equal(t, true, isCancelled)
		wg.Done()
	})

	assert.Nil(t, err)

	tb.Cancel()

	wg.Wait()
}

func TestExtend(t *testing.T) {

	tb := NewTimeBank()

	var wg sync.WaitGroup
	wg.Add(1)

	startTime := time.Now()

	err := tb.NewTask(3*time.Second, func(isCancelled bool) {
		wg.Done()
	})

	assert.Nil(t, err)

	tb.Extend(2 * time.Second)

	wg.Wait()

	duration := time.Now().Unix() - startTime.Unix()
	assert.Equal(t, int64(5), duration)
}
