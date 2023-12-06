package match

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/nats-io/nats.go"
)

type DispatchRequest struct {
	PlayerID               string `json:"player_id"`
	HighestNumberOfPlayers bool   `json:"highest_number_of_players"`
}

type Dispatcher interface {
	Start() error
	Close() error
	GetPendingCount() int64
	Dispatch(playerID string, highestNumberOfPlayers bool) error
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
	id := strings.Replace(d.m.Options().ID, "-", "", -1)
	queueName := fmt.Sprintf("match_dispatcher_%s", id)
	subject := fmt.Sprintf("match.dispatcher.%s", id)
	queue, err := d.m.QueueManager().AssertQueue(queueName, subject)
	if err != nil {
		fmt.Println(err)
		return err
	}

	d.queue = queue

	return d.queue.Subscribe(func(m *nats.Msg) {

		var req DispatchRequest

		json.Unmarshal(m.Data, &req)

		err := d.dispatch(&req)
		if err != nil {
			d.onFailure(err, req.PlayerID)
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

	err = d.queue.Destroy()
	if err != nil {
		return err
	}

	return nil
}

func (d *dispatcher) GetPendingCount() int64 {
	return atomic.LoadInt64(&d.pending)
}

// func (d *dispatcher) dispatch(playerID string) error {
func (d *dispatcher) dispatch(req *DispatchRequest) error {

	minAvailSeats := 1

	if d.m.IsLastTableStage() {
		// no limits
		minAvailSeats = -1
	}

	fmt.Printf("Dispatching player %s\n", req.PlayerID)

	// Find the table with the maximum number of players
	err := d.m.TableMap().DispatchPlayer(&TableCondition{
		HighestNumberOfPlayers: req.HighestNumberOfPlayers,
		MinAvailableSeats:      minAvailSeats,
	}, req.PlayerID)

	if err == ErrNotFoundAvailableTable {
		// No available table, so pushing to waiting room
		return d.m.WaitingRoom().Enter(req.PlayerID)
	}

	return err
}

func (d *dispatcher) Dispatch(playerID string, highestNumberOfPlayers bool) error {

	if d.queue == nil {
		return nil
	}

	atomic.AddInt64(&d.pending, int64(1))

	req := &DispatchRequest{
		PlayerID:               playerID,
		HighestNumberOfPlayers: highestNumberOfPlayers,
	}

	data, _ := json.Marshal(req)

	return d.queue.Publish(data)
}

func (d *dispatcher) OnFailure(fn func(error, string)) {
	d.onFailure = fn
}
