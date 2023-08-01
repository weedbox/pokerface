package match

import (
	"github.com/nats-io/nats.go"
)

type Queue interface {
	Assert() error
	Destroy() error
	Subject() string
	Publish(msg []byte) error
	Subscribe(fn func(*nats.Msg)) error
	Unsubscribe() error
}

type queue struct {
	qm      QueueManager
	name    string
	subject string
	sub     *nats.Subscription
}

func NewQueue(qm QueueManager, queueName string, subject string) Queue {
	return &queue{
		qm:      qm,
		name:    queueName,
		subject: subject,
	}
}

func (q *queue) Destroy() error {

	js, err := q.qm.Conn().JetStream()
	if err != nil {
		return err
	}

	_, err = js.StreamInfo(q.name)
	if err != nil && err == nats.ErrStreamNotFound {
		return nil
	}

	return js.DeleteStream(q.name)
}

func (q *queue) Assert() error {

	js, err := q.qm.Conn().JetStream()
	if err != nil {
		return err
	}

	_, err = js.StreamInfo(q.name)
	if err != nil && err != nats.ErrStreamNotFound {
		return err
	}

	// Exists already
	if err != nats.ErrStreamNotFound {
		return nil
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      q.name,
		Subjects:  []string{q.subject},
		Retention: nats.LimitsPolicy,
	})

	if err != nil {
		return err
	}

	return nil
}

func (q *queue) Subject() string {
	return q.subject
}

func (q *queue) Publish(msg []byte) error {

	js, err := q.qm.Conn().JetStream()
	if err != nil {
		return err
	}

	_, err = js.Publish(q.subject, msg)
	if err != nil {
		return err
	}

	return nil
}

func (q *queue) Subscribe(fn func(*nats.Msg)) error {

	js, err := q.qm.Conn().JetStream()
	if err != nil {
		return err
	}

	sub, err := js.Subscribe(q.subject, func(m *nats.Msg) {
		fn(m)
	})
	if err != nil {
		return err
	}

	q.sub = sub

	return nil
}

func (q *queue) Unsubscribe() error {
	return q.sub.Unsubscribe()
}
