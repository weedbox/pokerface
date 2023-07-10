package psae

type DummyRuntime struct {
}

func NewDummyRuntime() *DummyRuntime {
	return &DummyRuntime{}
}

func (dr *DummyRuntime) TableStateUpdated(p PSAE, ts *TableState) {
	// emit if table was updated
}

func (dr *DummyRuntime) TableBroken(p PSAE, ts *TableState) {
	// emit if table was updated
}

func (dr *DummyRuntime) PlayerJoined(p PSAE, player *Player) {
	// emit if player was joined
}

func (dr *DummyRuntime) PlayerDispatched(p PSAE, player *Player) {
	// emit if player was dispatched
}

func (dr *DummyRuntime) PlayerReleased(p PSAE, player *Player) {
	// emit if player was released
}

func (dr *DummyRuntime) Matched(p PSAE, match *Match) {
	// emit if matched
}

func (dr *DummyRuntime) WaitingRoomDrained(p PSAE, player *Player) {
	// emit if waiting room was drained
}

func (dr *DummyRuntime) WaitingRoomEntered(p PSAE, player *Player) {
	// emit if new player entered into waiting room
}

func (dr *DummyRuntime) WaitingRoomLeft(p PSAE, player *Player) {
	// emit if player left from waiting room
}

func (dr *DummyRuntime) WaitingRoomMatched(p PSAE, players []*Player) {
	// emit if satisfy condition for new table in the waiting room
}
