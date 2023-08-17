package match

import (
	"math"
)

type Runner interface {
	ShouldBeSplit(m Match, table *Table) bool
	DismissTable(m Match, table *Table) error
	DrainWaitingRoomPlayers(m Match, players []string) error
}

type NativeRunner struct {
}

func NewNativeRunner() *NativeRunner {
	return &NativeRunner{}
}

func (nr *NativeRunner) ShouldBeSplit(m Match, table *Table) bool {

	tableCount := m.TableMap().Count()

	// Final table
	if tableCount <= 1 {
		return false
	}

	totalPlayers := m.GetPlayerCount()

	//fmt.Printf("ShouldBeSplit (id=%s, table_count=%d, total_players=%d)\n", table.id, tableCount, totalPlayers)

	// The current number of tables is insufficient to accommodate all players.
	requiredTables := int(math.Ceil(float64(totalPlayers) / float64(m.Options().MaxSeats)))
	if int64(requiredTables) > tableCount {
		// It shouldn't break table to reduce the number of tables

		// table is freeze if only one player remains
		if table.GetPlayerCount() == 1 {
			// There are no other tables available to accommodate the remaining players from this table, so break this table
			return true
		}

		return false
	}

	if m.GetStatus() == MatchStatus_AfterRegDeadline {

		// The table is full, it should not be changed
		if table.GetPlayerCount() == m.Options().MaxSeats {
			return false
		}

		return true
	}

	// Condition 1: the number of players is less than or equal to minimum limit
	if table.GetPlayerCount() <= 3 {
		//fmt.Printf("table %s has players LESS THAN 3\n", table.ID())

		// Break table to release players
		return true
	}

	// Condition 2: the number of players is less than average
	if tableCount == 0 {
		return false
	}

	// Calculate the average number of players per table
	avg := totalPlayers / int64(tableCount)
	/*
		fmt.Printf("table count: %d, total players: %d, avg: %d, player in table: %d, no changes: %d\n",
			tableCount,
			totalPlayers,
			avg,
			table.GetPlayerCount(),
			table.noChanges,
		)
	*/
	if int64(table.GetPlayerCount()) < avg && table.noChanges >= m.Options().BreakThreshold {
		//fmt.Printf("table %s has players LASS THAN AVG\n", table.ID())
		// Break table to release players
		return true
	}

	return false
}

func (nr *NativeRunner) DismissTable(m Match, table *Table) error {

	players, err := table.GetPlayers()
	if err != nil {
		return err
	}

	for _, p := range players {
		//fmt.Printf("Releasing player %s from table %s\n", p, table.ID())
		err := m.Dispatcher().Dispatch(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (nr *NativeRunner) DrainWaitingRoomPlayers(m Match, players []string) error {

	// re-dispatch players who is drained from waiting room
	for _, id := range players {
		err := m.Dispatcher().Dispatch(id)
		if err != nil {
			return err
		}
	}

	return nil
}
