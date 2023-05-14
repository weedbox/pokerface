package waitgroup

func WaitReady(runtime *WaitGroupRuntime) bool {

	// Check if all players is ready
	for _, s := range runtime.States {
		if s.State == nil {
			return false
		}

		if s.State.(bool) == false {
			return false
		}
	}

	return true
}

func NewWaitReadyRuntime(players []int) *WaitGroupRuntime {

	states := make([]*WaitGroupPlayerState, 0, len(players))

	for _, seatSeq := range players {
		states = append(states, &WaitGroupPlayerState{
			Idx:   seatSeq,
			State: false,
		})
	}

	return NewWaitGroupRuntime(states)
}
