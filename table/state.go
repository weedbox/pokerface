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
	Status    string               `json:"status"`
	Options   *Options             `json:"options"`
	Players   map[int]*PlayerInfo  `json:"player"`
	GameState *pokerface.GameState `json:"game_state"`
}

func NewState() *State {
	return &State{
		Players: make(map[int]*PlayerInfo),
	}
}

func (s *State) GetJSON() []byte {

	data, _ := json.Marshal(s)

	return data
}

func (s *State) Clone() *State {

	// clone table state
	data, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	var state State
	json.Unmarshal(data, &state)

	return &state
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

func (s *State) GetPlayerByGameIdx(idx int) *PlayerInfo {
	for _, p := range s.Players {
		if p.GameIdx == idx {
			return p
		}
	}

	return nil
}

func (s *State) GetPlayerByID(playerID string) *PlayerInfo {
	for _, p := range s.Players {
		if p.ID == playerID {
			return p
		}
	}

	return nil
}

func (s *State) GetPlayerBySeatID(seatID int) *PlayerInfo {
	for _, p := range s.Players {
		if p.SeatID == seatID {
			return p
		}
	}

	return nil
}
