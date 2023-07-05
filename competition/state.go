package competition

type State struct {
	ID        string `json:"id"`
	GameType  string `json:"game_type"`
	StartTime int64  `json:"start_time"`
	State     string `json:"status"`
}

func NewState() *State {
	return &State{
		GameType: "standard",
		State:    "standby",
	}
}
