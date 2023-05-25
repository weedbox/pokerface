package task

type WaitPay struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Payload   int64  `json:"payload"`
}

func NewWaitPay(name string) *WaitPay {
	return &WaitPay{
		Name:    name,
		Payload: 0,
	}
}

func (wr *WaitPay) Instance() interface{} {
	return wr
}

func (wr *WaitPay) GetName() string {
	return wr.Name
}

func (wr *WaitPay) GetPayload() interface{} {
	return wr.Payload
}

func (wr *WaitPay) IsCompleted() bool {
	return wr.Completed
}

func (wr *WaitPay) Execute() bool {

	if wr.Payload > 0 {
		wr.Completed = true
	}

	return wr.Completed
}

func (wr *WaitPay) Pay(chips int64) {
	wr.Payload = chips
}
