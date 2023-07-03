package actor

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface/table"
)

func Test_NativeTableAdapter_Basic(t *testing.T) {

	backend := table.NewNativeBackend()
	opts := table.NewOptions()

	nt := table.NewTable(opts, table.WithBackend(backend))
	nt.SetAnte(10)
	nt.SetBlinds(0, 5, 10)

	// Initializing bots
	players := map[string]int64{
		"player_1": 10000,
		"player_2": 10000,
		"player_3": 10000,
	}

	// Preparing actors
	actors := make([]Actor, 0)
	for id, bankroll := range players {

		sid, _ := nt.Join(-1, &table.PlayerInfo{
			ID:       id,
			Bankroll: bankroll,
		})
		nt.Activate(sid)

		// Create new actor
		a := NewActor()

		// Initializing table engine adapter to communicate with native table
		ta := NewNativeTableAdapter(nt)
		a.SetAdapter(ta)

		// Initializing bot runner
		bot := NewBotRunner(id)
		a.SetRunner(bot)

		actors = append(actors, a)
	}

	// Setup state handler
	var wg sync.WaitGroup
	wg.Add(1)
	nt.OnStateUpdated(func(s *table.State) {

		// Update table state via adapter
		for _, a := range actors {
			a.GetTable().(*NativeTableAdapter).UpdateNativeState(s)
		}

		if s.Status == "playing" {
			if s.GameState.Status.CurrentEvent == "GameClosed" {
				t.Logf("GameClosed (id=%s, playable_players=%d)", s.GameState.GameID, nt.GetPlayablePlayerCount())
			}
		}

		if s.Status == "closed" {
			t.Log("TableClosed")
			assert.Less(t, nt.GetPlayablePlayerCount(), nt.GetState().Options.MinPlayers)
			wg.Done()
		}
	})

	// Not allow new player to join table
	nt.SetJoinable(false)

	assert.Nil(t, nt.Start())

	wg.Wait()
}
