package timebank

import (
	"errors"
	"time"
)

var (
	ErrInvalidParameters = errors.New("timebank: invalid parameters")
)

const (
	DefaultTimeout time.Duration = 15 * time.Second
)

type TimeBank struct {
	isRunning bool
	timer     *time.Timer
	due       time.Time
	callback  func(bool)
	closed    chan struct{}
}

func NewTimeBank() *TimeBank {

	// Initializing timer
	timer := time.NewTimer(DefaultTimeout)
	timer.Stop()

	tb := &TimeBank{
		isRunning: false,
		timer:     timer,
		closed:    make(chan struct{}),
	}

	return tb
}

func (tb *TimeBank) Cancel() {
	tb.closed <- struct{}{}
	tb.isRunning = false
	tb.timer.Stop()

	if tb.callback != nil {
		tb.callback(true)
	}
}

func (tb *TimeBank) NewTask(duration time.Duration, fn func(isCancelled bool)) error {

	if duration == time.Second*0 {
		return nil
	}

	if fn == nil {
		return ErrInvalidParameters
	}

	// Running already
	if tb.isRunning {
		tb.Cancel()
	}

	tb.due = time.Now().Add(duration)
	tb.timer.Reset(duration)
	tb.isRunning = true
	tb.callback = fn

	go func() {
		select {
		case <-tb.timer.C:
			tb.callback(false)
		case <-tb.closed:
		}
	}()

	return nil
}

func (tb *TimeBank) Extend(duration time.Duration) bool {

	// Time bank is not running
	if !tb.isRunning || tb.due.Before(time.Now()) {
		return false
	}

	// total = remain + extend
	total := tb.due.Sub(time.Now()) + duration
	tb.timer.Stop()

	// Update timer
	tb.due = tb.due.Add(duration)
	tb.timer.Reset(total)

	return true
}
