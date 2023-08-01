package match

type SeatChanges struct {
	Dealer int            `json:"dealer"` // The seat number of the dealer
	SB     int            `json:"sb"`     // The seat number of small blind
	BB     int            `json:"bb"`     // The seat number of big blind
	Seats  map[int]string `json:"seats"`  // Seat states (eg: left)
}

func NewSeatChanges() *SeatChanges {

	sc := &SeatChanges{
		Dealer: -1,
		SB:     -1,
		BB:     -1,
		Seats:  make(map[int]string),
	}

	return sc
}
