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
	tableEngine := pokertable.NewTableEngine()
	table, err := tableEngine.CreateTable(
		pokertable.TableSetting{
			ShortID:        "ABC123",
			Code:           "01",
			Name:           "table name",
			InvitationCode: "come_to_play",
			CompetitionMeta: pokertable.CompetitionMeta{
				ID: "competition id",
				Blind: pokertable.Blind{
					ID:              uuid.New().String(),
					Name:            "blind name",
					FinalBuyInLevel: 2,
					InitialLevel:    1,
					Levels: []pokertable.BlindLevel{
						{
							Level:    1,
							SB:       10,
							BB:       20,
							Ante:     0,
							Duration: 10,
						},
						{
							Level:    2,
							SB:       20,
							BB:       30,
							Ante:     0,
							Duration: 10,
						},
						{
							Level:    3,
							SB:       30,
							BB:       40,
							Ante:     0,
							Duration: 10,
						},
					},
				},
				MaxDuration:         60,
				Rule:                pokertable.CompetitionRule_Default,
				Mode:                pokertable.CompetitionMode_MTT,
				TableMaxSeatCount:   9,
				TableMinPlayerCount: 2,
				MinChipUnit:         10,
				ActionTime:          1,
			},
		},
	)
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
		for _, a := range actors {
			a.GetTable().UpdateTableState(table)
		}

		if table.State.Status == pokertable.TableStateStatus_TableGameSettled {
			if table.State.GameState.Status.CurrentEvent.Name == "GameClosed" {
				t.Log("GameClosed", table.State.GameState.GameID)

				if len(table.AlivePlayers()) == 1 {
					tableEngine.DeleteTable(table.ID)
					t.Log("Table deleted")
					wg.Done()
					return
				}
			}
		}
	})

	// Add player to table
	for _, p := range players {
		err := tableEngine.PlayerJoin(table.ID, p)
		assert.Nil(t, err)
	}

	// Start game
	err = tableEngine.StartTableGame(table.ID)
	assert.Nil(t, err)

	wg.Wait()
}
