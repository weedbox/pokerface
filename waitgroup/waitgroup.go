package waitgroup

type WaitGroupType int32

const (
	TypeReady WaitGroupType = iota
	TypePayAnte
)

type WaitGroupCheckFunc func(runtime *WaitGroupRuntime) bool

type WaitGroupPlayerState struct {
	Idx   int         `json:"idx"`
	State interface{} `json:"state"`
}

type WaitGroupRuntime struct {
	States []*WaitGroupPlayerState `json:"states"`
}

func NewWaitGroupRuntime(defStates []*WaitGroupPlayerState) *WaitGroupRuntime {
	return &WaitGroupRuntime{
		States: defStates,
	}
}

type WaitGroup struct {
	Type      WaitGroupType
	Runtime   *WaitGroupRuntime
	CheckFunc WaitGroupCheckFunc
}

func NewWaitGroup(t WaitGroupType, runtime *WaitGroupRuntime, checkFunc WaitGroupCheckFunc) *WaitGroup {
	return &WaitGroup{
		Type:      t,
		Runtime:   runtime,
		CheckFunc: checkFunc,
	}
}

func (wg *WaitGroup) IsCompleted() bool {
	return wg.CheckFunc(wg.Runtime)
}

func (wg *WaitGroup) GetStateByIdx(idx int) *WaitGroupPlayerState {

	for _, s := range wg.Runtime.States {
		if s.Idx == idx {
			return s
		}
	}

	return nil
}
