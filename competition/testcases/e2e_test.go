package competition

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface/actor"
	"github.com/weedbox/pokerface/competition"
	"github.com/weedbox/pokerface/table"
)

func Test_E2E(t *testing.T) {

	var tableEvents sync.Map
	var tables sync.Map

	initTable := func(tableID string) {

		// Load table instance to update states
		v, _ := tableEvents.Load(tableID)
		te := v.(chan *table.State)

		// New table so preparing handler
		go func(te chan *table.State) {
			var mu sync.Mutex
			for ts := range te {
				mu.Lock()

				// Loading all actors
				v, _ := tables.Load(tableID)
				actors := v.(*sync.Map)

				// Broadcast
				actors.Range(func(k interface{}, v interface{}) bool {
					a := v.(actor.Actor)
					err := a.GetTable().(*actor.NativeTableAdapter).UpdateNativeState(ts)
					if err != nil {
						t.Logf("Error: %v (player=%s, event=%s)", err, k.(string), ts.GameState.Status.CurrentEvent)
					}
					return true
				})
				mu.Unlock()
			}
		}(te)
	}

	assertTable := func(tableID string) chan *table.State {

		tables.LoadOrStore(tableID, &sync.Map{})

		// Load table instance to update states
		v, loaded := tableEvents.LoadOrStore(tableID, make(chan *table.State, 1024))
		te := v.(chan *table.State)

		if !loaded {
			initTable(tableID)
		}

		return te
	}

	closeTable := func(tableID string) {

		// Remove table
		v, _ := tableEvents.Load(tableID)
		te := v.(chan *table.State)
		close(te)

		tableEvents.Delete(tableID)
	}

	opts := competition.NewOptions()
	opts.TableAllocationPeriod = 1
	opts.MaxTables = -1
	opts.Table.Interval = 500

	var wg sync.WaitGroup
	wg.Add(1)

	closedCount := 0
	tb := competition.NewNativeTableBackend(table.NewNativeBackend())
	c := competition.NewCompetition(
		opts,
		competition.WithTableBackend(tb),
		competition.WithSeatReservedCallback(func(ts *table.State, seatID int, playerID string) {

			v, _ := tables.Load(ts.ID)
			actors := v.(*sync.Map)

			// Getting native table
			table := tb.(*competition.NativeTableBackend).GetTable(ts.ID)

			// Create new actor to join table
			a := actor.NewActor()

			// Initializing table engine adapter to communicate with table
			ta := actor.NewNativeTableAdapter(table)
			a.SetAdapter(ta)

			// Initializing bot runner
			bot := actor.NewBotRunner(playerID)
			a.SetRunner(bot)

			actors.Store(playerID, a)

			// Activate seats on specific table
			table.Activate(seatID)

			t.Logf("Player Seated (table=%s, seat=%d, player=%s)", ts.ID, seatID, playerID)
		}),
		competition.WithTableUpdatedCallback(func(ts *table.State) {

			te := assertTable(ts.ID)
			te <- ts

			if ts.Status == "playing" {

				//t.Log(ts.GameState.Status.CurrentEvent)
				if ts.GameState != nil && ts.GameState.Status.CurrentEvent == "GameClosed" {

					t.Logf("GameClosed (table=%s, id=%s, playable_players=%d)", ts.ID[:8], ts.GameState.GameID, len(ts.Players))
					/*
						for _, p := range ts.Players {
							t.Logf("================ (table=%s, seat=%d, player=%s, playable=%v)", ts.ID, p.SeatID, p.ID, p.Playable)
						}

						if playableCount == 1 {
							ts.PrintState()
						}
					*/
				}
			}

			if ts.Status == "closed" {
				closedCount++
				t.Logf("[%d] TableClosed (table=%s, players=%d)", closedCount, ts.ID[:8], len(ts.Players))
				closeTable(ts.ID)
			}
		}),
		competition.WithCompletedCallback(func(c competition.Competition) {
			t.Log("Completed")
			wg.Done()
		}),
	)
	defer c.Close()

	assert.Nil(t, c.Start())

	// Registering
	totalPlayer := 900
	for i := 0; i < totalPlayer; i++ {
		playerID := fmt.Sprintf("player_%d", i+1)
		assert.Nil(t, c.Register(playerID, 10000))
	}

	time.Sleep(time.Second)

	assert.Equal(t, totalPlayer, c.GetCompetitorCount())

	c.SetJoinable(false)

	wg.Wait()
}
