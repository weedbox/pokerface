package psae

type TestRuntimeOptions struct {
	TableStateUpdated  func(p PSAE, ts *TableState)
	TableBroken        func(p PSAE, ts *TableState)
	PlayerJoined       func(p PSAE, player *Player)
	PlayerDispatched   func(p PSAE, player *Player)
	PlayerReleased     func(p PSAE, player *Player)
	Matched            func(p PSAE, match *Match)
	WaitingRoomDrained func(p PSAE, player *Player)
	WaitingRoomEntered func(p PSAE, player *Player)
	WaitingRoomLeft    func(p PSAE, player *Player)
	WaitingRoomMatched func(p PSAE, players []*Player)
}

func NewTestRuntimeOptions() *TestRuntimeOptions {
	return &TestRuntimeOptions{
		TableStateUpdated:  func(p PSAE, ts *TableState) {},
		TableBroken:        func(p PSAE, ts *TableState) {},
		PlayerJoined:       func(p PSAE, player *Player) {},
		PlayerDispatched:   func(p PSAE, player *Player) {},
		PlayerReleased:     func(p PSAE, player *Player) {},
		Matched:            func(p PSAE, match *Match) {},
		WaitingRoomDrained: func(p PSAE, player *Player) {},
		WaitingRoomEntered: func(p PSAE, player *Player) {},
		WaitingRoomLeft:    func(p PSAE, player *Player) {},
		WaitingRoomMatched: func(p PSAE, players []*Player) {},
	}
}

type TestRuntime struct {
	opts *TestRuntimeOptions
}

func NewTestRuntime(opts *TestRuntimeOptions) *TestRuntime {
	return &TestRuntime{
		opts: opts,
	}
}

func (tr *TestRuntime) TableStateUpdated(p PSAE, ts *TableState) {
	// emit if table state was updated
	go tr.opts.TableStateUpdated(p, ts)
}

func (tr *TestRuntime) TableBroken(p PSAE, ts *TableState) {
	// emit if table was updated
	go tr.opts.TableBroken(p, ts)
}

func (tr *TestRuntime) PlayerJoined(p PSAE, player *Player) {
	// emit if player was joined
	go tr.opts.PlayerJoined(p, player)
}

func (tr *TestRuntime) PlayerDispatched(p PSAE, player *Player) {
	// emit if player was dispatched
	go tr.opts.PlayerDispatched(p, player)
}

func (tr *TestRuntime) PlayerReleased(p PSAE, player *Player) {
	// emit if player was released
	go tr.opts.PlayerReleased(p, player)
}

func (tr *TestRuntime) Matched(p PSAE, match *Match) {
	// emit if player was released
	go tr.opts.Matched(p, match)
}

func (tr *TestRuntime) WaitingRoomDrained(p PSAE, player *Player) {
	// emit if waiting room was trained
	go tr.opts.WaitingRoomDrained(p, player)
}

func (tr *TestRuntime) WaitingRoomEntered(p PSAE, player *Player) {
	// emit if new player entered into waiting room
	go tr.opts.WaitingRoomEntered(p, player)
}

func (tr *TestRuntime) WaitingRoomLeft(p PSAE, player *Player) {
	// emit if player left from waiting room
	go tr.opts.WaitingRoomLeft(p, player)
}

func (tr *TestRuntime) WaitingRoomMatched(p PSAE, players []*Player) {
	// emit if satisfy condition for new table in the waiting room
	go tr.opts.WaitingRoomMatched(p, players)
}
