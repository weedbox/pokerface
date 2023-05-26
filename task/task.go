package task

type TaskHandler func(ct *CustomizedTask) bool

type Task interface {
	GetType() string
	GetName() string
	GetPayload() interface{}
	IsCompleted() bool
	Execute() bool

	OnUpdated(fn func())
	OnCompleted(fn func())
}

type CustomizedTask struct {
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	Completed bool        `json:"completed"`
	Payload   interface{} `json:"payload"`

	handler     TaskHandler
	onUpdated   func()
	onCompleted func()
}

func NewTask(t string, name string, fn TaskHandler) Task {
	return &CustomizedTask{
		Type:    t,
		Name:    name,
		handler: fn,

		onUpdated:   func() {},
		onCompleted: func() {},
	}
}

func (ct *CustomizedTask) GetType() string {
	return ct.Type
}

func (ct *CustomizedTask) GetName() string {
	return ct.Name
}

func (ct *CustomizedTask) GetPayload() interface{} {
	return ct.Payload
}

func (ct *CustomizedTask) IsCompleted() bool {
	return ct.Completed
}

func (ct *CustomizedTask) Execute() bool {

	ct.Completed = ct.handler(ct)

	ct.onUpdated()

	if ct.Completed {
		ct.onCompleted()
	}

	return ct.Completed
}

func (ct *CustomizedTask) OnUpdated(fn func()) {
	ct.onUpdated = fn
}

func (ct *CustomizedTask) OnCompleted(fn func()) {
	ct.onCompleted = fn
}
