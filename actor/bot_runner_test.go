package actor

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	pokertable "github.com/weedbox/pokertable"
)

func TestActor_BotRunner_Humanize(t *testing.T) {

	// Initializing table
	manager := pokertable.NewManager()
	table, err := manager.CreateTable(pokertable.TableSetting{
		TableID: uuid.New().String(),
		Meta: pokertable.TableMeta{
			CompetitionID:       "1005c477-84b4-4d1b-9fca-3a6ad84e0fe7",
			Rule:                pokertable.CompetitionRule_Default,
			Mode:                pokertable.CompetitionMode_CT,
			MaxDuration:         3,
			TableMaxSeatCount:   9,
			TableMinPlayerCount: 2,
			MinChipUnit:         10,
			ActionTime:          10,
		},
	})
	assert.Nil(t, err)

	tableEngine, err := manager.GetTableEngine(table.ID)
	assert.Nil(t, err)

	// Initializing bot
	players := []pokertable.JoinPlayer{
		{PlayerID: "Jeffrey", RedeemChips: 3000},
		{PlayerID: "Chuck", RedeemChips: 3000},
		{PlayerID: "Fred", RedeemChips: 3000},
	}

	// Preparing actors
	actors := make([]Actor, 0)
	for _, p := range players {

		// Create new actor
		a := NewActor()

		// Initializing table engine adapter to communicate with table engine
		tc := NewTableEngineAdapter(tableEngine, table)
		a.SetAdapter(tc)

		// Initializing bot runner
		bot := NewBotRunner(p.PlayerID)
		bot.Humanized(true)
		a.SetRunner(bot)

		actors = append(actors, a)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Preparing table state updater
	tableEngine.OnTableUpdated(func(table *pokertable.Table) {

		// Update table state via adapter
		go func() {
			for _, a := range actors {
				a.GetTable().UpdateTableState(table)
			}
		}()

		if table.State.Status == pokertable.TableStateStatus_TableGameSettled {
			if table.State.GameState.Status.CurrentEvent == "GameClosed" {
				t.Log("GameClosed", table.State.GameState.GameID)

				if len(table.AlivePlayers()) == 1 {
					tableEngine.CloseTable()
					t.Log("Table deleted")
					wg.Done()
					return
				}
			}
		}
	})

	// Add player to table
	for _, p := range players {
		tableEngine.PlayerReserve(p)
		err := tableEngine.PlayerJoin(p.PlayerID)
		assert.Nil(t, err)
	}

	// Start game
	err = tableEngine.StartTableGame()
	assert.Nil(t, err)

	wg.Wait()
}
