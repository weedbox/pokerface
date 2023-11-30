package match

type Options struct {
	ID                string `json:"id"`
	WaitingPeriod     int    `json:"waiting_period"`
	MinInitialPlayers int    `json:"min_initial_player"`
	MaxSeats          int    `json:"max_seats"`
	MaxTables         int    `json:"max_tables"`
	Joinable          bool   `json:"joinable"`
	BreakThreshold    int    `json:"break_threshold"`
}

func NewOptions(id string) *Options {
	return &Options{
		ID:                id,
		WaitingPeriod:     10, // 10 seconds
		MinInitialPlayers: 4,
		MaxSeats:          9,
		MaxTables:         -1, // Unlimit
		Joinable:          true,
		BreakThreshold:    3,
	}
}
