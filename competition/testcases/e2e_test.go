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

	var tables sync.Map

	opts := competition.NewOptions()
	opts.TableAllocationPeriod = 1
	opts.MaxTables = -1
	opts.Table.Interval = 500

	var wg sync.WaitGroup
	wg.Add(1)

	tb := competition.NewNativeTableBackend(table.NewNativeBackend())
	c := competition.NewCompetition(
		opts,
		competition.WithTableBackend(tb),
		competition.WithPlayerJoinedCallback(func(ts *table.State, seatID int, playerID string) {

			v, _ := tables.LoadOrStore(ts.ID, &sync.Map{})
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

			t.Logf("Joined (table=%s, seat=%d, player=%s)", ts.ID, seatID, playerID)
		}),
		competition.WithTableUpdatedCallback(func(ts *table.State) {

			v, ok := tables.Load(ts.ID)
			assert.True(t, ok)
			actors := v.(*sync.Map)

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

			// Update table state via adapter
			go actors.Range(func(k interface{}, v interface{}) bool {
				a := v.(actor.Actor)
				err := a.GetTable().(*actor.NativeTableAdapter).UpdateNativeState(ts)
				if err != nil {
					t.Logf("Error: %v (player=%s, event=%s)", err, k.(string), ts.GameState.Status.CurrentEvent)
				}
				return true
			})

			if ts.Status == "closed" {
				time.Sleep(time.Second)
				t.Log("TableClosed")
				assert.Less(t, len(ts.Players), ts.Options.MinPlayers)
				//wg.Done()
			}
		}),
	)
	defer c.Close()

	assert.Nil(t, c.Start())

	// Registering
	totalPlayer := 18
	for i := 0; i < totalPlayer; i++ {
		playerID := fmt.Sprintf("player_%d", i+1)
		assert.Nil(t, c.Register(playerID, 10000))
	}

	time.Sleep(time.Second)

	c.SetJoinable(false)

	wg.Wait()
}
