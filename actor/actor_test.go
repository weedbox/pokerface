package actor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pokertable "github.com/weedbox/pokertable"
)

func TestActor_Basic(t *testing.T) {

	// Initializing table
	tableEngine := pokertable.NewTableEngine(uint32(logrus.DebugLevel))
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
							Level:        1,
							SBChips:      10,
							BBChips:      20,
							AnteChips:    0,
							DurationMins: 10,
						},
						{
							Level:        2,
							SBChips:      20,
							BBChips:      30,
							AnteChips:    0,
							DurationMins: 10,
						},
						{
							Level:        3,
							SBChips:      30,
							BBChips:      40,
							AnteChips:    0,
							DurationMins: 10,
						},
					},
				},
				MaxDurationMins:      60,
				Rule:                 pokertable.CompetitionRule_Default,
				Mode:                 pokertable.CompetitionMode_MTT,
				TableMaxSeatCount:    9,
				TableMinPlayingCount: 2,
				MinChipsUnit:         10,
			},
		},
	)
	assert.Nil(t, err)

	// Initializing bot
	players := []pokertable.JoinPlayer{
		{PlayerID: "Jeffrey", RedeemChips: 150},
		{PlayerID: "Chuck", RedeemChips: 150},
		{PlayerID: "Fred", RedeemChips: 150},
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
		a.SetRunner(bot)

		actors = append(actors, a)
	}

	// Preparing table state updater
	tableEngine.OnTableUpdated(func(table *pokertable.Table) {

		for _, a := range actors {
			a.UpdateTableState(table)
		}
	})

	// Add player to table
	for _, p := range players {
		err := tableEngine.PlayerJoin(table.ID, p)
		assert.Nil(t, err)
	}

	// Start game
	err = tableEngine.StartGame(table.ID)
	assert.Nil(t, err)
}
