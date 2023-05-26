package task

type WaitReady struct {
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Completed bool         `json:"completed"`
	Payload   map[int]bool `json:"payload"`

	onUpdated   func()
	onCompleted func()
}

func NewWaitReady(name string) *WaitReady {
	return &WaitReady{
		Type:    "ready",
		Name:    name,
		Payload: make(map[int]bool),

		onUpdated:   func() {},
		onCompleted: func() {},
	}
}

func (wr *WaitReady) GetType() string {
	return wr.Type
}

func (wr *WaitReady) GetName() string {
	return wr.Name
}

func (wr *WaitReady) GetPayload() interface{} {
	return wr.Payload
}

func (wr *WaitReady) IsCompleted() bool {
	return wr.Completed
}

func (wr *WaitReady) Execute() bool {

	for _, isReady := range wr.Payload {
		if !isReady {
			wr.onUpdated()
			return false
		}
	}

	wr.Completed = true
	wr.onUpdated()
	wr.onCompleted()

	return true
}

func (wr *WaitReady) PreparePlayerStates(playerCount int) {

	for i := 0; i < playerCount; i++ {
		wr.Payload[i] = false
	}
}

func (wr *WaitReady) Ready(playerIdx int) {

	if len(wr.Payload) <= playerIdx {
		return
	}

	wr.Payload[playerIdx] = true
}

func (wr *WaitReady) OnUpdated(fn func()) {
	wr.onUpdated = fn
}

func (wr *WaitReady) OnCompleted(fn func()) {
	wr.onCompleted = fn
}
