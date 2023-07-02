package table

import (
	"encoding/json"
	"fmt"

	"github.com/weedbox/pokerface"
)

type State struct {
	ID        string               `json:"id"`
	GameType  string               `json:"game_type"`
	StartTime int64                `json:"start_time"`
	EndTime   int64                `json:"end_time"`
	Players   map[int]*PlayerInfo  `json:"player"`
	Status    string               `json:"status"`
	GameState *pokerface.GameState `json:"game_state"`
}

func NewState() *State {
	return &State{
		Players: make(map[int]*PlayerInfo),
	}
}

func (s *State) PrintState() error {

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func (s *State) ResetPositions() {
	for _, p := range s.Players {
		p.Positions = make([]string, 0)
	}
}
