package table

import (
	"errors"
	"math/rand"
)

var (
	ErrNotFoundSeat                = errors.New("seat_manager: not found seat")
	ErrNoAvailableSeat             = errors.New("seat_manager: no available seat")
	ErrNotAvailable                = errors.New("seat_manager: not available")
	ErrInvalidSeat                 = errors.New("seat_manager: invalid seat")
	ErrInsufficientNumberOfPlayers = errors.New("seat_manager: insufficient number of players")
	ErrEmptySeat                   = errors.New("seat_manager: empty seat")
)

type Seat struct {
	ID         int
	Player     *PlayerInfo
	IsActive   bool
	IsReserved bool
}

type SeatManager struct {
	max         int
	seats       map[int]*Seat
	playerCount int

	dealer *Seat
	sb     *Seat
	bb     *Seat
}

func NewSeatManager(max int) *SeatManager {

	sm := &SeatManager{
		max:   max,
		seats: make(map[int]*Seat),
	}

	// Initializing seats
	for i := 0; i < max; i++ {
		sm.seats[i] = &Seat{
			ID:         i,
			Player:     nil,
			IsActive:   true,
			IsReserved: false,
		}
	}

	return sm
}

func (sm *SeatManager) renewSeatStatus() error {

	origSeats := sm.GetNormalizeSeats(sm.dealer.ID)
	seats := origSeats

	//	if sm.playerCount == 2 {
	if sm.GetPlayableSeatCount() == 2 {
		// dealer is SB as well
		sm.sb = sm.dealer
	} else {
		// Find SB based on current dealer
		seats = seats[1:]
		sb, idx := sm.findActivePlayer(seats)
		sm.sb = sb
		seats = seats[idx:]
	}

	// Find BB based on current SB
	seats = seats[1:]
	bb, idx := sm.findActivePlayer(seats)
	sm.bb = bb
	seats = seats[idx:]

	// Deactivate seats between dealer and BB
	for _, s := range origSeats {
		if s == sm.bb {
			break
		}

		if s.Player == nil {
			s.IsActive = false
		}
	}

	// Activate the rest of seats
	seats = seats[1:]
	for _, s := range seats {
		s.IsActive = true
	}

	return nil
}

func (sm *SeatManager) join(seatID int, p *PlayerInfo) (int, error) {

	s := sm.GetSeat(seatID)
	if s.Player != nil {
		return -1, ErrNotAvailable
	}

	s.IsReserved = true
	s.Player = p
	sm.playerCount++

	return s.ID, nil
}

func (sm *SeatManager) leave(seatID int) error {

	s := sm.GetSeat(seatID)
	if s.Player == nil {
		return ErrEmptySeat
	}

	s.Player = nil
	s.IsReserved = false
	sm.playerCount--

	return nil
}

func (sm *SeatManager) findActivePlayer(seats []*Seat) (*Seat, int) {

	for i, s := range seats {

		// Ignore seat which is not activated and empty
		if !s.IsActive || s.IsReserved || s.Player == nil {
			continue
		}

		// Found
		return s, i
	}

	return nil, -1
}

func (sm *SeatManager) NextDealer() *Seat {

	//	if sm.playerCount < 2 {
	if sm.GetPlayableSeatCount() < 2 {
		return nil
	}

	var seats []*Seat
	if sm.dealer == nil {
		// from the first seat
		seats = sm.GetNormalizeSeats(0)
	} else {
		// From the current dealer
		seats = sm.GetNormalizeSeats(sm.dealer.ID)
		seats = seats[1:]
	}

	// Find the next dealer
	dealer, _ := sm.findActivePlayer(seats)

	// Found
	if dealer != nil {

		// Activate seats between old dealer and new dealer
		for _, s := range seats {
			if s == dealer {
				break
			}

			s.IsActive = true
		}

		sm.dealer = dealer
		return dealer
	}

	// Not found the next dealer, because all of player has been left except new players who is inactive
	for _, s := range seats {
		s.IsActive = true
	}

	// Try again. It should get a new dealer as long as more than one players out there
	sm.dealer, _ = sm.findActivePlayer(seats)

	return sm.dealer
}

func (sm *SeatManager) GetDealer() *Seat {
	return sm.dealer
}

func (sm *SeatManager) GetSmallBlind() *Seat {
	return sm.sb
}

func (sm *SeatManager) GetBigBlind() *Seat {
	return sm.bb
}

func (sm *SeatManager) GetPlayerCount() int {
	return sm.playerCount
}

func (sm *SeatManager) GetSeat(id int) *Seat {

	if s, ok := sm.seats[id]; ok {
		return s
	}

	return nil
}

func (sm *SeatManager) GetNormalizeSeats(startID int) []*Seat {

	// Getting player list that specific seat should be the first element of it
	cur := startID

	seats := make([]*Seat, 0)
	for i := 0; i < sm.max; i++ {

		if s, ok := sm.seats[cur]; ok {
			seats = append(seats, s)
		}

		// next player
		cur++

		// The end of seat list
		if cur == sm.max {
			cur = 0
		}
	}

	return seats
}

func (sm *SeatManager) GetAvailableSeats() ([]int, []int) {

	seats := make([]int, 0)
	alternateSeats := make([]int, 0)

	for _, s := range sm.seats {

		if s.IsReserved || s.Player != nil {
			continue
		}

		if s.IsActive {
			seats = append(seats, s.ID)
		} else {
			alternateSeats = append(alternateSeats, s.ID)
		}
	}

	return seats, alternateSeats
}

func (sm *SeatManager) GetSeats() []*Seat {

	seats := make([]*Seat, 0)
	for i := 0; i < sm.max; i++ {
		s := sm.seats[i]
		seats = append(seats, s)
	}

	return seats
}

func (sm *SeatManager) GetActiveSeats() []*Seat {

	seats := make([]*Seat, 0)
	for i := 0; i < sm.max; i++ {
		s := sm.seats[i]
		if s.IsActive {
			seats = append(seats, s)
		}
	}

	return seats
}

func (sm *SeatManager) GetPlayableSeats() []*Seat {

	origSeats := sm.GetNormalizeSeats(sm.dealer.ID)

	seats := make([]*Seat, 0)
	for _, s := range origSeats {
		if !s.IsReserved && s.IsActive && s.Player != nil {
			seats = append(seats, s)
		}
	}

	return seats
}

func (sm *SeatManager) GetPlayableSeatCount() int {

	count := 0
	for i := 0; i < sm.max; i++ {
		s := sm.seats[i]
		if s.IsActive && !s.IsReserved && s.Player != nil {
			count++
		}
	}

	return count
}

func (sm *SeatManager) Activate(seatID int) error {

	seat := sm.GetSeat(seatID)
	if seat == nil {
		return ErrNotFoundSeat
	}

	seat.IsReserved = false

	return nil
}

func (sm *SeatManager) Reserve(seatID int) error {

	seat := sm.GetSeat(seatID)
	if seat == nil {
		return ErrNotFoundSeat
	}

	seat.IsReserved = true

	return nil
}

func (sm *SeatManager) Join(seatID int, p *PlayerInfo) (int, error) {

	if seatID >= sm.max {
		return -1, ErrInvalidSeat
	}

	// Specific seat
	if seatID > -1 {
		return sm.join(seatID, p)
	}

	// Getting available seats
	s, as := sm.GetAvailableSeats()
	if len(s) == 0 && len(as) == 0 {
		return -1, ErrNoAvailableSeat
	}

	// Select a seat from list randomly
	if len(s) > 0 {
		if len(s) == 1 {
			return sm.join(s[0], p)
		}

		return sm.join(s[rand.Intn(len(s)-1)], p)
	}

	// Worst-case: Select a seat from alternate seats betwwen dealer and BB
	if len(as) == 1 {
		return sm.join(as[0], p)
	}

	return sm.join(as[rand.Intn(len(as)-1)], p)
}

func (sm *SeatManager) Leave(seatID int) error {
	return sm.leave(seatID)
}

func (sm *SeatManager) Next() error {

	if sm.NextDealer() == nil {
		return ErrInsufficientNumberOfPlayers
	}

	return sm.renewSeatStatus()
}
