package actor

import (
	"encoding/json"
	"time"

	"github.com/weedbox/pokerface"
	"github.com/weedbox/pokerface/table"
	"github.com/weedbox/pokertable"
)

type NativeTableAdapter struct {
	actor Actor
	table table.Table
	state *table.State
}

func NewNativeTableAdapter(t table.Table) *NativeTableAdapter {

	return &NativeTableAdapter{
		table: t,
	}
}

func (nta *NativeTableAdapter) SetActor(a Actor) {
	nta.actor = a
}

func (nta *NativeTableAdapter) UpdateTableState(t *pokertable.Table) error {
	return nta.actor.UpdateTableState(t)
}

func (nta *NativeTableAdapter) UpdateNativeState(state *table.State) error {

	// Clone to get a individual table structure
	data := state.GetJSON()

	var s table.State
	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		return err
	}

	state = &s

	nta.state = state

	// Convert native table state to standard format
	t := pokertable.Table{
		ID: state.ID,
		Meta: pokertable.TableMeta{
			ActionTime: state.Options.ActionTime,
		},
		State: &pokertable.TableState{
			GameState:         state.GameState,
			PlayerStates:      make([]*pokertable.TablePlayerState, 0),
			GamePlayerIndexes: make([]int, 0),
		},
	}

	switch state.Status {
	case "idle":
		t.State.Status = pokertable.TableStateStatus_TableGameStandby
	case "preparing":
		t.State.Status = pokertable.TableStateStatus_TableGameStandby
	case "playing":
		t.State.Status = pokertable.TableStateStatus_TableGamePlaying
	case "closed":
		t.State.Status = pokertable.TableStateStatus_TableClosed
	}

	if state.GameState != nil {
		t.State.GamePlayerIndexes = make([]int, len(state.GameState.Players))
		//fmt.Println("Update players")

		count := 0
		for _, p := range state.Players {
			//fmt.Printf("<<<<=========== player=%s, gameIdx=%d, psIdx=%d\n", p.ID, p.GameIdx, count)
			t.State.PlayerStates = append(t.State.PlayerStates, &pokertable.TablePlayerState{
				PlayerID:  p.ID,
				Seat:      p.SeatID,
				Positions: p.Positions,
				Bankroll:  p.Bankroll,
			})

			if p.GameIdx != -1 {
				/*
					if p.GameIdx >= len(state.GameState.Players) {
						fmt.Println("=======", state.Status, state.GameState.Status.CurrentEvent)
						for _, gp := range state.GameState.Players {
							fmt.Println(gp.Idx)
						}
						for _, p := range state.Players {
							fmt.Printf("<<<<===== player=%s, gameIdx=%d, p.SeatID=%d, bankroll=%d, playable=%v\n", p.ID, p.GameIdx, p.SeatID, p.Bankroll, p.Playable)
						}

						state.PrintState()
					}
				*/
				t.State.GamePlayerIndexes[p.GameIdx] = count
			}

			count++
		}
	}

	return nta.UpdateTableState(&t)
}

func (nta *NativeTableAdapter) GetGameState() *pokerface.GameState {
	return nta.state.GameState
}

func (nta *NativeTableAdapter) GetGamePlayerIndex(playerID string) int {
	for _, p := range nta.state.Players {
		if p.ID == playerID {
			return p.GameIdx
		}
	}

	return -1
}

func (nta *NativeTableAdapter) Pass(playerID string) error {
	return nta.table.Pass(playerID)
}

func (nta *NativeTableAdapter) Ready(playerID string) error {
	return nta.table.Ready(playerID)
}

func (nta *NativeTableAdapter) Pay(playerID string, chips int64) error {
	return nta.table.Pay(playerID, chips)
}

func (nta *NativeTableAdapter) Check(playerID string) error {
	return nta.table.Check(playerID)
}

func (nta *NativeTableAdapter) Bet(playerID string, chips int64) error {
	return nta.table.Bet(playerID, chips)
}

func (nta *NativeTableAdapter) Call(playerID string) error {
	return nta.table.Call(playerID)
}

func (nta *NativeTableAdapter) Fold(playerID string) error {
	return nta.table.Fold(playerID)
}

func (nta *NativeTableAdapter) Allin(playerID string) error {
	return nta.table.Allin(playerID)
}

func (nta *NativeTableAdapter) Raise(playerID string, chipLevel int64) error {
	return nta.table.Raise(playerID, chipLevel)
}

func (nta *NativeTableAdapter) ExtendTime(playerID string, duration time.Duration) error {
	//TODO: need to be implemented
	return nil
}
