package match

import (
	"errors"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type QueueManager interface {
	Connect() error
	Close() error
	Conn() *nats.Conn
	AssertQueue(queueName string, subject string) (Queue, error)
}

type NativeQueueManager struct {
	server *server.Server
	nc     *nats.Conn
}

func NewNativeQueueManager() QueueManager {

	nqm := &NativeQueueManager{}

	return nqm
}

func (nqm *NativeQueueManager) Connect() error {

	opts := &server.Options{
		JetStream: true,
	}

	s, err := server.NewServer(opts)
	if err != nil {
		return err
	}

	nqm.server = s

	go s.Start()

	if !s.ReadyForConnections(4 * time.Second) {
		return errors.New("not ready for connection")
	}

	nc, err := nats.Connect(s.ClientURL())
	if err != nil {
		return err
	}

	nqm.nc = nc

	return nil
}

func (nqm *NativeQueueManager) Close() error {

	nqm.nc.Close()
	nqm.server.Shutdown()
	nqm.server.WaitForShutdown()

	return nil
}

func (nqm *NativeQueueManager) Conn() *nats.Conn {
	return nqm.nc
}

func (nqm *NativeQueueManager) AssertQueue(queueName string, subject string) (Queue, error) {

	q := NewQueue(nqm, queueName, subject)

	err := q.Destroy()
	if err != nil {
		return nil, err
	}

	err = q.Assert()
	if err != nil {
		return nil, err
	}

	return q, nil
}
