package psae

type Backend struct {
	AllocateTable func() (*TableState, error)
	JoinTable     func(tableID string, players []*Player) error
	BrokeTable    func(tableID string) error
}

func NewBackend() *Backend {
	return &Backend{
		AllocateTable: func() (*TableState, error) { return nil, nil },
		JoinTable:     func(tableID string, players []*Player) error { return nil },
		BrokeTable:    func(tableID string) error { return nil },
	}
}
