package psae

import "fmt"

var tableCounter = 0
var playerCounter = 0

func NewTestTableState(playerCount int) *TableState {
	tableCounter++
	symbol := fmt.Sprintf("table %d", tableCounter)
	ts := &TableState{
		ID:             symbol,
		Players:        make(map[string]*Player),
		Status:         TableStatus_Ready,
		TotalSeats:     9,
		AvailableSeats: 9,
		Statistics: &TableStatistics{
			NoChanges: 0,
		},
	}

	for i := 0; i < playerCount; i++ {
		p := NewTestPlayer()
		ts.Players[p.ID] = p
	}

	ts.AvailableSeats = ts.TotalSeats - playerCount

	return ts
}

func NewTestPlayer() *Player {
	playerCounter++
	symbol := fmt.Sprintf("player %d", playerCounter)
	return &Player{
		ID:   symbol,
		Name: symbol,
	}
}
