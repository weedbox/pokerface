package match

import (
	"math"

	"github.com/weedbox/pokerface/match/psae"
)

type MatchRuntime struct {
}

func NewMatchRuntime() *MatchRuntime {
	return &MatchRuntime{}
}

func (mr *MatchRuntime) TableStateUpdated(p psae.PSAE, ts *psae.TableState) {

	if p.GetStatus() == psae.GameStatus_AfterRegistrationDeadline {

		count, err := p.SeatMap().GetTableCount()
		if err != nil {
			return
		}

		// Final table
		if count == 1 {
			// Do nothing
			return
		}

		totalPlayers, err := p.SeatMap().GetTotalPlayers()
		if err != nil {
			return
		}

		// too many tables
		if math.Ceil(float64(totalPlayers)/float64(p.Game().MaxPlayersPerTable)) < float64(count) {

			// table is full, it should not be changed
			if len(ts.Players) == p.Game().MaxPlayersPerTable {
				// Do nothing
				return
			}

			// break table to release players
			ts.Status = psae.TableStatus_Broken
			err := p.BreakTable(ts.ID)
			if err != nil {
				return
			}
		}

		return
	}

	// Condition 1: number of players are less than or equal to minimum limit
	if len(ts.Players) <= 3 {
		//fmt.Printf("table %s has players LESS THAN 3\n", ts.ID)
		// breaking table
		ts.Status = psae.TableStatus_Broken
		err := p.BreakTable(ts.ID)
		if err != nil {
			return
		}

		return
	}

	// Condition 2: number of players are less than average
	tableCount, err := p.SeatMap().GetTableCount()
	if err != nil {
		return
	}

	if tableCount == 0 {
		return
	}

	// Calculate average number of players
	totalPlayers, err := p.SeatMap().GetTotalPlayers()
	if err != nil {
		return
	}

	avg := totalPlayers / int64(tableCount)
	//fmt.Printf("table count: %d, total players: %d, avg: %d, player in table: %d , no changes: %d\n", tableCount, totalPlayers, avg, len(ts.Players), ts.Statistics.NoChanges)
	if int64(len(ts.Players)) < avg && ts.Statistics.NoChanges >= 10 {

		//fmt.Printf("table %s has players LASS THAN AVG\n", ts.ID)

		// breaking table
		ts.Status = psae.TableStatus_Broken
		err := p.BreakTable(ts.ID)
		if err != nil {
			return
		}
	}
}

func (mr *MatchRuntime) TableBroken(p psae.PSAE, ts *psae.TableState) {

	// Release players from the table
	for _, player := range ts.Players {
		p.ReleasePlayer(player)
	}
}

func (mr *MatchRuntime) PlayerJoined(p psae.PSAE, player *psae.Player) {
}

func (mr *MatchRuntime) PlayerDispatched(p psae.PSAE, player *psae.Player) {

	minAvailSeats := 1
	if p.IsLastTableStage() {
		minAvailSeats = -1
	}

	// Find a table which has highest number of players
	target, err := p.SeatMap().FindAvailableTable(&psae.TableCondition{
		HighestNumberOfPlayers: true,
		MinAvailableSeats:      minAvailSeats,
	})

	if err != nil {
		return
	}

	// Not found
	if target == nil {

		//fmt.Printf("Stay in waiting room (player=%s)\n", player.ID)

		// Stay in waiting room
		err = p.EnterWaitingRoom(player)
		if err != nil {
			return
		}

		return
	}

	// Join target table
	err = p.JoinTable(target.ID, []*psae.Player{
		player,
	})

	if err != nil {
		return
	}
}

func (mr *MatchRuntime) PlayerReleased(p psae.PSAE, player *psae.Player) {

	minAvailSeats := 1
	if p.IsLastTableStage() {
		minAvailSeats = -1
	}

	// Find a table which has lowest number of players
	target, err := p.SeatMap().FindAvailableTable(&psae.TableCondition{
		HighestNumberOfPlayers: false,
		MinAvailableSeats:      minAvailSeats,
	})

	//fmt.Printf("Found new table %s (player_count=%d, avail=%d) for player (%s)\n", target.ID, len(target.Players), target.AvailableSeats, player.ID)

	if err != nil {
		return
	}

	// Not found
	if target == nil {
		// Re-dispatch just like new player
		err = p.DispatchPlayer(player)
		if err != nil {
			return
		}

		return
	}

	// Join target table
	err = p.JoinTable(target.ID, []*psae.Player{
		player,
	})

	if err != nil {
		return
	}
}

func (mr *MatchRuntime) Matched(p psae.PSAE, m *psae.Match) {

	// Preparing a new table
	ts, err := p.AllocateTable()
	if err != nil {
		return
	}

	ts.Status = psae.TableStatus_Suspend
	err = p.AssertTableState(ts)
	if err != nil {
		return
	}

	// For matched players
	err = p.JoinTable(ts.ID, m.Players)
	if err != nil {
		return
	}

	// Put players into the table
	ts.AvailableSeats = ts.TotalSeats - len(m.Players)
	for _, player := range m.Players {
		ts.Players[player.ID] = player
	}

	// Registering in seat map
	ts.Status = psae.TableStatus_Ready
	_, err = p.UpdateTableState(ts)
	if err != nil {
		return
	}
}

func (mr *MatchRuntime) WaitingRoomDrained(p psae.PSAE, player *psae.Player) {

	// Re-dispatch just like new player
	err := p.DispatchPlayer(player)
	if err != nil {
		return
	}
}

func (mr *MatchRuntime) WaitingRoomEntered(p psae.PSAE, player *psae.Player) {
}

func (mr *MatchRuntime) WaitingRoomLeft(p psae.PSAE, player *psae.Player) {
}

func (mr *MatchRuntime) WaitingRoomMatched(p psae.PSAE, players []*psae.Player) {

	err := p.MatchPlayers(players)
	if err != nil {
		return
	}
}
