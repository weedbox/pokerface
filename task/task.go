package task

type TaskHandler func(ct *CustomizedTask) bool

type Task interface {
	Instance() interface{}
	GetName() string
	GetPayload() interface{}
	IsCompleted() bool
	Execute() bool
}

type CustomizedTask struct {
	Name      string      `json:"name"`
	Completed bool        `json:"completed"`
	Payload   interface{} `json:"payload"`

	handler TaskHandler
}

func NewTask(name string, fn TaskHandler) Task {
	return &CustomizedTask{
		Name:    name,
		handler: fn,
	}
}

func (ct *CustomizedTask) Instance() interface{} {
	return ct
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

	return ct.Completed
}
