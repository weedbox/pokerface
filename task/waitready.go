package task

type WaitReady struct {
	Name      string       `json:"name"`
	Completed bool         `json:"completed"`
	Payload   map[int]bool `json:"payload"`
}

func NewWaitReady(name string) *WaitReady {
	return &WaitReady{
		Name:    name,
		Payload: make(map[int]bool),
	}
}

func (wr *WaitReady) Instance() interface{} {
	return wr
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
			return false
		}
	}

	wr.Completed = true

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
