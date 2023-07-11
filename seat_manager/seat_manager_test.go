package seat_manager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPlayerInfo struct {
	ID        string
	Positions []string
}

func Test_SeatManager_Join(t *testing.T) {

	sm := NewSeatManager(9)

	for i := 1; i < 10; i++ {
		seatID, err := sm.Join(i-1, &TestPlayerInfo{
			ID:        fmt.Sprintf("Player %d", i),
			Positions: make([]string, 0),
		})
		assert.Nil(t, err)
		assert.Nil(t, sm.Activate(seatID))
	}

	assert.Equal(t, 9, sm.GetPlayerCount())

	s, as := sm.GetAvailableSeats()
	assert.Equal(t, 0, len(s))
	assert.Equal(t, 0, len(as))

	// Attempt to specify seat which is unavailable
	_, err := sm.Join(1, &TestPlayerInfo{
		ID:        fmt.Sprintf("Player %d", 1),
		Positions: make([]string, 0),
	})
	assert.Equal(t, ErrNotAvailable, err)
}

func Test_SeatManager_Join_Random(t *testing.T) {

	sm := NewSeatManager(9)

	for i := 1; i < 10; i++ {
		seatID, err := sm.Join(-1, &TestPlayerInfo{
			ID:        fmt.Sprintf("Player %d", i),
			Positions: make([]string, 0),
		})

		assert.Nil(t, err)
		assert.Nil(t, sm.Activate(seatID))
	}

	assert.Equal(t, 9, sm.GetPlayerCount())

	seats := sm.GetSeats()
	for _, s := range seats {
		assert.NotNil(t, s.Player)
	}
}

func Test_SeatManager_ReservedSeat(t *testing.T) {

	sm := NewSeatManager(9)

	for i := 1; i < 10; i++ {
		_, err := sm.Join(i-1, &TestPlayerInfo{
			ID:        fmt.Sprintf("Player %d", i),
			Positions: make([]string, 0),
		})
		assert.Nil(t, err)
	}

	assert.Equal(t, 9, sm.GetPlayerCount())
	assert.Equal(t, 0, sm.GetPlayableSeatCount())
	assert.Equal(t, ErrInsufficientNumberOfPlayers, sm.Next())

	for i := 1; i < 10; i++ {
		assert.Nil(t, sm.Activate(i-1))
		assert.Equal(t, i, sm.GetPlayableSeatCount())
	}
}

func Test_SeatManager_Next(t *testing.T) {

	sm := NewSeatManager(9)

	for i := 1; i < 10; i++ {
		seatID, err := sm.Join(-1, &TestPlayerInfo{
			ID:        fmt.Sprintf("Player %d", i),
			Positions: make([]string, 0),
		})

		assert.Nil(t, err)
		assert.Nil(t, sm.Activate(seatID))
	}

	// Cycle
	for i := 0; i < sm.GetPlayerCount()+1; i++ {

		dealerIdx := i
		if dealerIdx >= sm.GetPlayerCount() {
			dealerIdx = dealerIdx - sm.GetPlayerCount()
		}

		sbIdx := i + 1
		if sbIdx >= sm.GetPlayerCount() {
			sbIdx = sbIdx - sm.GetPlayerCount()
		}

		bbIdx := i + 2
		if bbIdx >= sm.GetPlayerCount() {
			bbIdx = bbIdx - sm.GetPlayerCount()
		}

		// First game
		assert.Nil(t, sm.Next())

		seats := sm.GetSeats()
		assert.Equal(t, sm.GetDealer(), seats[dealerIdx])
		assert.Equal(t, seats[dealerIdx], sm.dealer)
		assert.Equal(t, sm.GetSmallBlind(), seats[sbIdx])
		assert.Equal(t, seats[sbIdx], sm.sb)
		assert.Equal(t, sm.GetBigBlind(), seats[bbIdx])
		assert.Equal(t, seats[bbIdx], sm.bb)
	}
}

func Test_SeatManager_Next_TwoPlayer(t *testing.T) {

	sm := NewSeatManager(9)

	for i := 1; i < 3; i++ {
		seatID, err := sm.Join(i-1, &TestPlayerInfo{
			ID:        fmt.Sprintf("Player %d", i),
			Positions: make([]string, 0),
		})

		assert.Nil(t, err)
		assert.Nil(t, sm.Activate(seatID))
	}

	// Cycle
	for i := 0; i < sm.GetPlayerCount()+1; i++ {

		dealerIdx := i
		if dealerIdx >= sm.GetPlayerCount() {
			dealerIdx = dealerIdx - sm.GetPlayerCount()
		}

		bbIdx := i + 1
		if bbIdx >= sm.GetPlayerCount() {
			bbIdx = bbIdx - sm.GetPlayerCount()
		}

		assert.Nil(t, sm.Next())

		seats := sm.GetSeats()

		// Dealer is the small blind
		assert.Equal(t, sm.GetDealer(), seats[dealerIdx])
		assert.Equal(t, seats[dealerIdx], sm.dealer)
		assert.Equal(t, sm.GetSmallBlind(), seats[dealerIdx])
		assert.Equal(t, seats[dealerIdx], sm.sb)
		assert.Equal(t, sm.GetBigBlind(), seats[bbIdx])
		assert.Equal(t, seats[bbIdx], sm.bb)
	}
}

func Test_SeatManager_GetAvailableSeats(t *testing.T) {

	sm := NewSeatManager(9)
	s, as := sm.GetAvailableSeats()

	assert.Equal(t, 0, sm.GetPlayerCount())
	assert.Equal(t, 9, len(s))
	assert.Equal(t, 0, len(as))

	sm.Join(0, &TestPlayerInfo{
		ID:        "Player 1",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(0))

	sm.Join(2, &TestPlayerInfo{
		ID:        "Player 2",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(2))

	sm.Join(4, &TestPlayerInfo{
		ID:        "Player 3",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(4))

	assert.Equal(t, 3, sm.GetPlayerCount())

	// Initial positions
	assert.Nil(t, sm.Next())

	seats := sm.GetSeats()
	assert.Equal(t, sm.GetDealer(), seats[0])
	assert.Equal(t, sm.GetSmallBlind(), seats[2])
	assert.Equal(t, sm.GetBigBlind(), seats[4])

	// Getting available seats again
	s, as = sm.GetAvailableSeats()

	assert.Equal(t, 4, len(s))  // SeatID=5,6,7,8
	assert.Equal(t, 2, len(as)) // SeatID=1,3 (Between Dealer and SB)
}

func Test_SeatManager_AlternateSeats(t *testing.T) {

	sm := NewSeatManager(9)
	s, as := sm.GetAvailableSeats()

	assert.Equal(t, 0, sm.GetPlayerCount())
	assert.Equal(t, 9, len(s))
	assert.Equal(t, 0, len(as))

	sm.Join(0, &TestPlayerInfo{
		ID:        "Player 1",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(0))

	sm.Join(2, &TestPlayerInfo{
		ID:        "Player 2",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(2))

	sm.Join(4, &TestPlayerInfo{
		ID:        "Player 3",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(4))

	assert.Equal(t, 3, sm.GetPlayerCount())

	// Initial positions
	assert.Nil(t, sm.Next())

	// Getting available seats again
	s, as = sm.GetAvailableSeats()

	// Seat between Dealer and BB
	sm.Join(1, &TestPlayerInfo{
		ID:        "Player 4",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(1))
	sm.Join(3, &TestPlayerInfo{
		ID:        "Player 5",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(3))

	// Seat after BB
	sm.Join(5, &TestPlayerInfo{
		ID:        "Player 6",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(5))

	// Move position of game
	assert.Nil(t, sm.Next())

	seats := sm.GetSeats()

	assert.Equal(t, sm.GetDealer(), seats[2])
	assert.Equal(t, sm.GetSmallBlind(), seats[4])
	assert.Equal(t, sm.GetBigBlind(), seats[5])

	// Keep moving (Dealer=4,5,0,1)
	assert.Nil(t, sm.Next())
	assert.Nil(t, sm.Next())
	assert.Nil(t, sm.Next())
	assert.Nil(t, sm.Next())

	// Now Seat 1 and 3 can play
	assert.Equal(t, sm.GetDealer(), seats[1])
	assert.Equal(t, sm.GetSmallBlind(), seats[2])
	assert.Equal(t, sm.GetBigBlind(), seats[3])
}

func Test_SeatManager_AlternateSeats_Rejoin(t *testing.T) {

	sm := NewSeatManager(9)
	s, as := sm.GetAvailableSeats()

	assert.Equal(t, 0, sm.GetPlayerCount())
	assert.Equal(t, 9, len(s))
	assert.Equal(t, 0, len(as))

	sm.Join(0, &TestPlayerInfo{
		ID:        "Player 1",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(0))

	sm.Join(2, &TestPlayerInfo{
		ID:        "Player 2",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(2))

	sm.Join(4, &TestPlayerInfo{
		ID:        "Player 3",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(4))

	assert.Equal(t, 3, sm.GetPlayerCount())

	// initial positions
	assert.Nil(t, sm.Next())

	seats := sm.GetSeats()
	assert.Equal(t, sm.GetDealer(), seats[0])
	assert.Equal(t, sm.GetSmallBlind(), seats[2])
	assert.Equal(t, sm.GetBigBlind(), seats[4])

	// Player leaves
	sm.Leave(4)

	// Dealer is 2
	assert.Nil(t, sm.Next())

	seats = sm.GetSeats()
	assert.Equal(t, sm.GetDealer(), seats[2])
	assert.Equal(t, sm.GetSmallBlind(), seats[2])
	assert.Equal(t, sm.GetBigBlind(), seats[0])

	// re-join
	sm.Join(4, &TestPlayerInfo{
		ID:        "Player 3",
		Positions: make([]string, 0),
	})
	assert.Nil(t, sm.Activate(4))

	// Dealer is 0 (skip the new player)
	assert.Nil(t, sm.Next())

	seats = sm.GetSeats()
	assert.Equal(t, sm.GetDealer(), seats[0])
	assert.Equal(t, sm.GetSmallBlind(), seats[2])
	assert.Equal(t, sm.GetBigBlind(), seats[4])

	// Dealer is 2 (new player can player now)
	assert.Nil(t, sm.Next())

	// Now Seat 4 can play
	assert.Equal(t, sm.GetDealer(), seats[2])
	assert.Equal(t, sm.GetSmallBlind(), seats[4])
	assert.Equal(t, sm.GetBigBlind(), seats[0])
}
