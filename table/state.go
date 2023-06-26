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

func (s *State) PrintState() error {

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

type PlayerInfo struct {
	ID        string   `json:"id"`
	SeatID    int      `json:"seat_id"`
	GameIdx   int      `json:"game_idx"`
	Positions []string `json:"positions"`
	Bankroll  int64    `json:"bankroll"`
}

func NewState() *State {
	return &State{
		Players: make(map[int]*PlayerInfo),
	}
}

func (pi *PlayerInfo) CheckPosition(pos string) bool {

	for _, position := range pi.Positions {
		if pos == position {
			return true
		}
	}

	return false
}

func (pi *PlayerInfo) Assign(pos string) {

	for _, position := range pi.Positions {
		if pos == position {
			return
		}
	}

	pi.Positions = append(pi.Positions, pos)
}
