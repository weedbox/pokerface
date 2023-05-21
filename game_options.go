package main

import "github.com/cfsghost/pokerface/combination"

type GameOptions struct {
	Ante                   int64                     `json:"ante"`
	Blind                  BlindSetting              `json:"blind"`
	Limit                  string                    `json:"limit"`
	HoleCardsCount         int                       `json:"hole_cards_count"`
	RequiredHoleCardsCount int                       `json:"required_hole_cards_count"`
	CombinationPowers      []combination.Combination `json:"combination_powers"`
	Deck                   []string                  `json:"deck"`
	BurnCount              int                       `json:"burn_count"`
	Players                []*PlayerSetting          `json:"players"`
}

type BlindSetting struct {
	Dealer int64 `json:"dealer"`
	SB     int64 `json:"sb"`
	BB     int64 `json:"bb"`
}

type PlayerSetting struct {
	Bankroll  int64    `json:"bankroll"`
	Positions []string `json:"positions"`
}

func NewStardardGameOptions() *GameOptions {
	return &GameOptions{
		Ante: 0,
		Blind: BlindSetting{
			Dealer: 0,
			SB:     5,
			BB:     10,
		},
		Limit:                  "no",
		HoleCardsCount:         2,
		RequiredHoleCardsCount: 0,
		CombinationPowers:      combination.CombinationPowerStandard,
		//Deck
		BurnCount: 1,
		Players:   make([]*PlayerSetting, 0),
	}
}
func NewShortDeckGameOptions() *GameOptions {

	opts := NewStardardGameOptions()
	opts.CombinationPowers = combination.CombinationPowerShortDeck
	opts.RequiredHoleCardsCount = 2

	return opts
}
