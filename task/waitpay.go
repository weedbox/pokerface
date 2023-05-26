package task

type WaitPay struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Payload   int64  `json:"payload"`

	onUpdated   func()
	onCompleted func()
}

func NewWaitPay(name string) *WaitPay {
	return &WaitPay{
		Type:    "pay",
		Name:    name,
		Payload: 0,

		onUpdated:   func() {},
		onCompleted: func() {},
	}
}

func (wp *WaitPay) GetType() string {
	return wp.Type
}

func (wp *WaitPay) GetName() string {
	return wp.Name
}

func (wp *WaitPay) GetPayload() interface{} {
	return wp.Payload
}

func (wp *WaitPay) IsCompleted() bool {
	return wp.Completed
}

func (wp *WaitPay) Execute() bool {

	if wp.Payload > 0 {
		wp.Completed = true
	}

	wp.onUpdated()

	if wp.Completed {
		wp.onCompleted()
	}

	return wp.Completed
}

func (wp *WaitPay) Pay(chips int64) {
	wp.Payload = chips
}

func (wp *WaitPay) OnCompleted(fn func()) {
	wp.onCompleted = fn
}

func (wp *WaitPay) OnUpdated(fn func()) {
	wp.onUpdated = fn
}
