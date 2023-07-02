package table

import "github.com/weedbox/pokerface"

type Options struct {
	GameType       string                 `json:"game_type"`
	InitialPlayers int                    `json:"initial_players"`
	MinPlayers     int                    `json:"min_players"`
	MaxSeats       int                    `json:"max_seats"`
	MaxGames       int                    `json:"max_games"`
	Duration       int                    `json:"duration"`
	Interval       int                    `json:"interval"`
	ActionTime     int                    `json:"action_time"`
	Ante           int64                  `json:"ante"`
	Blind          pokerface.BlindSetting `json:"blind"`
}

func NewOptions() *Options {
	return &Options{
		GameType:       "standard",
		InitialPlayers: 2,
		MinPlayers:     2,
		MaxSeats:       9,
		MaxGames:       0,       // unlimit by default
		Duration:       60 * 60, // one hour
		Interval:       0,       // 0 secs by default
		ActionTime:     10,      // 10 secs
		Ante:           0,
		Blind: pokerface.BlindSetting{
			Dealer: 0,
			SB:     5,
			BB:     10,
		},
	}
}
