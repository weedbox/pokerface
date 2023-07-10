package psae

type MemoryPlayerQueue struct {
	queue chan *Player
}

func NewMemoryPlayerQueue() *MemoryPlayerQueue {
	return &MemoryPlayerQueue{
		queue: make(chan *Player, 10000),
	}
}

func (mq *MemoryPlayerQueue) Publish(p *Player) error {
	mq.queue <- p
	return nil
}

func (mq *MemoryPlayerQueue) Subscribe() (chan *Player, error) {
	return mq.queue, nil
}

func (mq *MemoryPlayerQueue) Close() error {
	close(mq.queue)
	return nil
}
