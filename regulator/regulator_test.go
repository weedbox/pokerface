package regulator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegulator(t *testing.T) {

	tableCounter := 0

	r := NewRegulator(
		WithRequestTableFn(func(players []string) (string, error) {
			tableCounter++
			t.Log("Request to create table", tableCounter)

			for _, player := range players {
				t.Log("  Player", player, "joined table")
			}

			return fmt.Sprintf("table_%d", tableCounter), nil
		}),

		WithAssignPlayersFn(func(tableID string, players []string) error {
			t.Log("Request to assign players to table", tableID)

			for _, player := range players {
				t.Log("  Assigned player", player, "to table")
			}

			return nil
		}),
	)

	totalPlayers := 0

	for i := 0; i < 9; i++ {
		totalPlayers++
		r.AddPlayers([]string{fmt.Sprintf("player_%d", totalPlayers)})
	}

	assert.Equal(t, 9, r.GetPlayerCount())
	assert.Equal(t, 0, r.GetTableCount())

	r.SetStatus(CompetitionStatus_Normal)

	assert.Equal(t, 1, r.GetTableCount())

	for i := 0; i < 3; i++ {
		totalPlayers++
		r.AddPlayers([]string{fmt.Sprintf("player_%d", totalPlayers)})
	}

	assert.Equal(t, 12, r.GetPlayerCount())
	assert.Equal(t, 2, r.GetTableCount())

	// table 2 has 3 players now, but it needs 6
	assert.Equal(t, 3, r.GetTable("table_2").Required)

	// Table still has 9 when first hand is over
	releaseCount, players, err := r.SyncState("table_1", 9)
	assert.Nil(t, err)

	// water level should be 6 for 12 players, so 3 players should be released
	assert.Equal(t, 3, releaseCount)

	// No new players should be put on table 1
	assert.Len(t, players, 0)

	// release players
	releasedPlayers := []string{
		"player_1",
		"player_2",
		"player_3",
	}
	err = r.ReleasePlayers("table_1", releasedPlayers)
	assert.Nil(t, err)

	// Table 2 should have 6 players now
	assert.Equal(t, 6, r.GetTable("table_2").PlayerCount)
	assert.Equal(t, 0, r.GetTable("table_2").Required)
}

func TestRegulator_91Problem(t *testing.T) {

	// Prepare tables
	tables := make(map[string][]string)
	for i := 0; i < 10; i++ {
		tableID := fmt.Sprintf("table_%d", i+1)
		tables[tableID] = []string{}
	}

	tableCounter := 0

	r := NewRegulator(
		WithRequestTableFn(func(players []string) (string, error) {
			tableCounter++
			tableID := fmt.Sprintf("table_%d", tableCounter)

			t.Log("Request to create table: ", tableID)

			for _, player := range players {
				t.Log("  Player", player, "joined table")
				tables[tableID] = append(tables[tableID], player)
			}

			return tableID, nil
		}),

		WithAssignPlayersFn(func(tableID string, players []string) error {
			t.Log("Request to assign players to table", tableID)

			for _, player := range players {
				t.Log("  Assigned player", player, "to table")
			}

			return nil
		}),
	)

	totalPlayers := 0

	// Prepare 90 players
	for i := 0; i < 90; i++ {
		totalPlayers++
		r.AddPlayers([]string{fmt.Sprintf("player_%d", totalPlayers)})
	}

	assert.Equal(t, 90, r.GetPlayerCount())
	assert.Equal(t, 0, r.GetTableCount())

	r.SetStatus(CompetitionStatus_Normal)

	// Add one more player
	totalPlayers++
	r.AddPlayers([]string{fmt.Sprintf("player_%d", totalPlayers)})

	// table 11 has 1 players now, but it still needs 7 players
	assert.Equal(t, 7, r.GetTable("table_11").Required)

	totalRequired := 0
	for i := 0; i < 10; i++ {

		tableID := fmt.Sprintf("table_%d", i+1)

		// Each table still has 9 when first hand is over
		releaseCount, players, err := r.SyncState(tableID, 9)
		assert.Nil(t, err)

		// No new players should be put on old table
		assert.Len(t, players, 0)

		totalRequired += releaseCount

		// Attempt to release players
		var releasedPlayers []string
		table := tables[tableID]
		for n := 0; n < releaseCount; n++ {

			// Pick one player to release
			player := table[0]
			table = table[1:]

			releasedPlayers = append(releasedPlayers, player)
		}

		tables[tableID] = table

		err = r.ReleasePlayers(tableID, releasedPlayers)
	}

	assert.Equal(t, 7, totalRequired)
	assert.Equal(t, 0, r.GetTable("table_11").Required)
}

func TestRegulator_AfterRegDeadline(t *testing.T) {

	tables := make(map[string][]string)

	tableCounter := 0

	r := NewRegulator(
		WithRequestTableFn(func(players []string) (string, error) {
			tableCounter++
			tableID := fmt.Sprintf("table_%d", tableCounter)
			tables[tableID] = []string{}

			t.Log("Request to create table: ", tableID)

			for _, player := range players {
				t.Log("  Player", player, "joined table")
				tables[tableID] = append(tables[tableID], player)
			}

			return tableID, nil
		}),

		WithAssignPlayersFn(func(tableID string, players []string) error {
			t.Log("Request to assign players to table", tableID)

			for _, player := range players {
				t.Log("  Assigned player", player, "to table")
				tables[tableID] = append(tables[tableID], player)
			}

			return nil
		}),
	)

	totalPlayers := 0

	// Prepare 27 players
	for i := 0; i < 27; i++ {
		totalPlayers++
		r.AddPlayers([]string{fmt.Sprintf("player_%d", totalPlayers)})
	}

	assert.Equal(t, 27, r.GetPlayerCount())
	assert.Equal(t, 0, r.GetTableCount())

	r.SetStatus(CompetitionStatus_Normal)
	r.SetStatus(CompetitionStatus_AfterRegDeadline)

	for i := 0; i < 3; i++ {

		t.Log("hand", i+1)

		for tableID, players := range tables {

			// 3 players are out
			players = players[3:]

			// Each table still has 9 when first hand is over
			releaseCount, newPlayers, err := r.SyncState(tableID, len(players))
			assert.Nil(t, err)

			t.Logf("Table %s: %d players, %d new players, should release %d players", tableID, len(players), len(newPlayers), releaseCount)

			// Attempt to release players
			var releasedPlayers []string
			for n := 0; n < releaseCount; n++ {

				// Pick one player to release
				player := players[0]
				players = players[1:]

				releasedPlayers = append(releasedPlayers, player)
			}

			err = r.ReleasePlayers(tableID, releasedPlayers)
			assert.Nil(t, err)

			// It should break this table
			if len(players) == 0 {
				t.Log("Break table:", tableID)
				delete(tables, tableID)
				continue
			}

			// Attempt to allocate for new players
			players = append(players, newPlayers...)

			// Update local table information
			tables[tableID] = players
		}

		if i == 0 {
			assert.Equal(t, 18, r.GetPlayerCount())
			assert.Equal(t, 2, r.GetTableCount())
		} else if i == 1 {
			assert.Equal(t, 12, r.GetPlayerCount())
			assert.Equal(t, 2, r.GetTableCount())
		} else if i == 2 {
			assert.Equal(t, 6, r.GetPlayerCount())
			assert.Equal(t, 1, r.GetTableCount())
		}
	}
}
