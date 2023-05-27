package task

type PaymentRequest struct {
	Required int64         `json:"required"`
	Players  map[int]int64 `json:"players"`
}

type WaitPay struct {
	Type      string         `json:"type"`
	Name      string         `json:"name"`
	Completed bool           `json:"completed"`
	Payload   PaymentRequest `json:"payload"`

	onUpdated   func()
	onCompleted func()
}

func NewWaitPay(name string, requiredChips int64) *WaitPay {
	return &WaitPay{
		Type: "pay",
		Name: name,
		Payload: PaymentRequest{
			Required: requiredChips,
			Players:  make(map[int]int64),
		},

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

	for _, chips := range wp.Payload.Players {

		if chips == 0 {
			wp.onUpdated()
			return false
		}
	}

	wp.Completed = true
	wp.onUpdated()
	wp.onCompleted()

	return wp.Completed
}

func (wp *WaitPay) PrepareStates(players []int) {

	for _, idx := range players {
		wp.Payload.Players[idx] = 0
	}
}

func (wp *WaitPay) Pay(playerIdx int, chips int64) {

	if chips == 0 {
		return
	}

	wp.Payload.Players[playerIdx] = chips
}

func (wp *WaitPay) OnCompleted(fn func()) {
	wp.onCompleted = fn
}

func (wp *WaitPay) OnUpdated(fn func()) {
	wp.onUpdated = fn
}
