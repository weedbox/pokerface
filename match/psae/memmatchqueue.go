package psae

type MemoryMatchQueue struct {
	queue chan *Match
}

func NewMemoryMatchQueue() *MemoryMatchQueue {
	return &MemoryMatchQueue{
		queue: make(chan *Match, 10000),
	}
}

func (mq *MemoryMatchQueue) Publish(p *Match) error {
	mq.queue <- p
	return nil
}

func (mq *MemoryMatchQueue) Subscribe() (chan *Match, error) {
	return mq.queue, nil
}

func (mq *MemoryMatchQueue) Close() error {
	close(mq.queue)
	return nil
}
