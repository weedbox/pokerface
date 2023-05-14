package waitgroup

func WaitPayAnte(runtime *WaitGroupRuntime) bool {

	// Check if all players paid ante ready
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

func NewWaitPayAnteRuntime(players []int) *WaitGroupRuntime {

	states := make([]*WaitGroupPlayerState, 0, len(players))

	for _, seatSeq := range players {
		states = append(states, &WaitGroupPlayerState{
			Idx:   seatSeq,
			State: false,
		})
	}

	return NewWaitGroupRuntime(states)
}
