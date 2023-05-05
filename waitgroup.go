package main

type PlayerWaitState struct {
	Idx        int   `json:"idx"`
	State      int32 `json:"state"`
	IsAnswered bool  `json:"is_answered"`
}

type WaitGroup struct {
	States []*PlayerWaitState
}

func NewWaitGroup(defStates []*PlayerWaitState) *WaitGroup {
	return &WaitGroup{
		States: defStates,
	}
}
