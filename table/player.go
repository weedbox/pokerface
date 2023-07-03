package table

type PlayerInfo struct {
	ID        string   `json:"id"`
	SeatID    int      `json:"seat_id"`
	GameIdx   int      `json:"game_idx"`
	Positions []string `json:"positions"`
	Playable  bool     `json:"playable"`
	Bankroll  int64    `json:"bankroll"`
}

func (pi *PlayerInfo) CheckPosition(pos string) bool {

	for _, position := range pi.Positions {
		if pos == position {
			return true
		}
	}

	return false
}

func (pi *PlayerInfo) Assign(pos string) {

	for _, position := range pi.Positions {
		if pos == position {
			return
		}
	}

	pi.Positions = append(pi.Positions, pos)
}
