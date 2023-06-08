package actor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	pokertable "github.com/weedbox/pokertable"
	"github.com/weedbox/pokertable/model"
	pokertableModel "github.com/weedbox/pokertable/model"
	"github.com/weedbox/pokertable/util"
)

func TestActor_Basic(t *testing.T) {

	// Preparing table
	gameEngine := pokertable.NewGameEngine()
	tableEngine := pokertable.NewTableEngine(gameEngine)
	table, _ := tableEngine.CreateTable(
		pokertableModel.TableSetting{
			ShortID:        "ABC123",
			Code:           "01",
			Name:           "table name",
			InvitationCode: "come_to_play",
			CompetitionMeta: model.CompetitionMeta{
				ID: "competition id",
				Blind: model.Blind{
					ID:              uuid.New().String(),
					Name:            "blind name",
					FinalBuyInLevel: 2,
					InitialLevel:    1,
					Levels: []model.BlindLevel{
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
				Rule:                 util.CompetitionRule_Default,
				Mode:                 util.CompetitionMode_MTT,
				TableMaxSeatCount:    9,
				TableMinPlayingCount: 2,
				MinChipsUnit:         10,
			},
		},
	)

	// Initializing bot
	players := []pokertableModel.JoinPlayer{
		{PlayerID: "Jeffrey", RedeemChips: 150},
		{PlayerID: "Chuck", RedeemChips: 150},
		{PlayerID: "Fred", RedeemChips: 150},
	}

	actors := make([]Actor, 0)
	for _, p := range players {
		// Create new actor
		a := NewActor()

		// Initializing bot runner
		bot := NewBotRunner(p.PlayerID)
		a.SetRunner(bot)

		// Initializing table engine adapter to communicate with table engine
		tc := NewTableEngineAdapter(tableEngine, table)
		a.SetAdapter(tc)

		actors = append(actors, a)
	}

	// Start game
	_, err := tableEngine.StartGame(table)
	assert.Nil(t, err)
}
