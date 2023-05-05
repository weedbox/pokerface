package main

type GameOptions struct {
	Ante                   int64         `json:"ante"`
	Blind                  BlindSetting  `json:"blind"`
	Limit                  string        `json:"limit"`
	HoleCardsCount         int           `json:"hole_cards_count"`
	RequiredHoleCardsCount int           `json:"required_hole_cards_count"`
	CombinationPowers      []Combination `json:"combination_powers"`
	Deck                   []string      `json:"deck"`
	BurnCount              int           `json:"burn_count"`
}

type BlindSetting struct {
	Dealer int64 `json:"dealer"`
	SB     int64 `json:"sb"`
	BB     int64 `json:"bb"`
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
		CombinationPowers:      CombinationPowerStandard,
		//Deck
		BurnCount: 1,
	}
}
func NewShortDeckGameOptions() *GameOptions {

	opts := NewStardardGameOptions()
	opts.CombinationPowers = CombinationPowerShortDeck
	opts.RequiredHoleCardsCount = 2

	return opts
}
