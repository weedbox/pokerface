package match

import (
	"fmt"
	"sync/atomic"

	"github.com/nats-io/nats.go"
)

type Dispatcher interface {
	Start() error
	Close() error
	GetPendingCount() int64
	Dispatch(playerID string) error
	OnFailure(func(err error, playerID string))
}

type dispatcher struct {
	m         Match
	queue     Queue
	onFailure func(err error, playerID string)
	pending   int64
}

func NewDispatcher(m Match) Dispatcher {
	return &dispatcher{
		m:         m,
		onFailure: func(error, string) {},
	}
}

func (d *dispatcher) Start() error {

	queue, err := d.m.QueueManager().AssertQueue("match_dispatcher", "match.dispatcher")
	if err != nil {
		fmt.Println(err)
		return err
	}

	d.queue = queue

	return d.queue.Subscribe(func(m *nats.Msg) {

		playerID := string(m.Data)
		err := d.dispatch(playerID)
		if err != nil {
			d.onFailure(err, playerID)
		}
		atomic.AddInt64(&d.pending, -int64(1))

		m.Ack()
	})
}

func (d *dispatcher) Close() error {

	if d.queue == nil {
		return nil
	}

	err := d.queue.Unsubscribe()
	if err != nil {
		return err
	}

	return nil
}

func (d *dispatcher) GetPendingCount() int64 {
	return atomic.LoadInt64(&d.pending)
}

func (d *dispatcher) dispatch(playerID string) error {

	minAvailSeats := 1

	if d.m.IsLastTableStage() {
		// no limits
		minAvailSeats = -1
	}

	//fmt.Printf("Dispatching player %s\n", playerID)

	// Find the table with the maximum number of players
	err := d.m.TableMap().DispatchPlayer(&TableCondition{
		HighestNumberOfPlayers: true,
		MinAvailableSeats:      minAvailSeats,
	}, playerID)

	if err == ErrNotFoundAvailableTable {
		// No available table, so pushing to waiting room
		return d.m.WaitingRoom().Enter(playerID)
	}

	return err
}

func (d *dispatcher) Dispatch(playerID string) error {

	if d.queue == nil {
		return nil
	}

	atomic.AddInt64(&d.pending, int64(1))
	return d.queue.Publish([]byte(playerID))
}

func (d *dispatcher) OnFailure(fn func(error, string)) {
	d.onFailure = fn
}
