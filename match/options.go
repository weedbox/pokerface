package match

type Options struct {
	WaitingPeriod     int  `json:"waiting_period"`
	MinInitialPlayers int  `json:"min_initial_player"`
	MaxSeats          int  `json:"max_seats"`
	MaxTables         int  `json:"max_tables"`
	Joinable          bool `json:"joinable"`
}

func NewOptions() *Options {
	return &Options{
		WaitingPeriod:     10, // 10 seconds
		MinInitialPlayers: 4,
		MaxSeats:          9,
		MaxTables:         -1, // Unlimit
		Joinable:          true,
	}
}
