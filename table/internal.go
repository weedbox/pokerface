package table

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/weedbox/pokerface"
)

func (t *table) delay(interval int, fn func() error) error {

	var err error
	var wg sync.WaitGroup
	wg.Add(1)

	t.tb.NewTask(time.Duration(interval)*time.Second, func(isCancelled bool) {

		defer wg.Done()

		if isCancelled {
			return
		}

		err = fn()
	})

	wg.Wait()

	return err
}

func (t *table) setupPosition() error {

	if t.inPosition {
		return nil
	}

	// Calculating positions for players
	err := t.sm.Next()
	if err != nil {
		return err
	}

	// Updating seat and position information for players
	t.ts.ResetPositions()

	seats := t.sm.GetSeats()
	for _, s := range seats {

		if s.Player == nil {
			continue
		}

		p := t.GetPlayerByID(s.Player.ID)

		positions := make([]string, 0)

		if s == t.sm.dealer {
			positions = append(positions, "dealer")
		}

		if s == t.sm.sb {
			positions = append(positions, "sb")
		} else if s == t.sm.bb {
			positions = append(positions, "bb")
		}

		p.Positions = positions
	}

	t.inPosition = true

	return nil
}

func (t *table) updatePlayerStates(ts *State) error {

	if ts.GameState == nil {
		return nil
	}

	if ts.GameState.Status.CurrentEvent != "GameClosed" {
		return nil
	}

	// Updating player states with settlement
	for _, rs := range ts.GameState.Result.Players {

		p := t.GetPlayerByGameIdx(rs.Idx)
		if p == nil {
			continue
		}

		p.Bankroll = rs.Final

		// Reserve the seat because player is unplayable
		if p.Bankroll == 0 {
			t.sm.Reserve(p.SeatID)
		}
	}

	return nil
}

func (t *table) updateStates(gs *pokerface.GameState) error {

	t.ts.GameState = gs

	// clone table state
	data, err := json.Marshal(t.ts)
	if err != nil {
		return err
	}

	var state State
	json.Unmarshal(data, &state)

	t.updatePlayerStates(&state)

	go t.onStateUpdated(&state)

	if t.ts.Status == "idle" {
		//fmt.Println("Attempt to start the next game")
		t.delay(t.options.Interval, func() error {
			return t.nextGame(t.options.Interval)
		})
	}

	return nil
}

func (t *table) run(delay int) error {

	t.ts.Status = "preparing"

	if t.options.MaxGames > 0 && t.options.MaxGames == t.gameCount {
		return ErrMaxGamesExceeded
	}

	err := t.setupPosition()
	if err != nil {
		return err
	}

	// Check remaining time
	if time.Now().Unix() >= t.ts.EndTime {
		// Times up!
		return ErrTimesUp
	}

	// Check the number of player
	if len(t.sm.GetPlayableSeats()) < t.options.MinPlayers {
		return ErrGameConditionsNotMet
	}

	return t.delay(delay, func() error {
		// Starting a new game
		return t.startGame()
	})
}

func (t *table) nextGame(delay int) error {

	if t.isPaused {
		return nil
	}

	err := t.run(delay)

	switch err {
	case ErrMaxGamesExceeded:
		fallthrough
	case ErrTimesUp:
		t.ts.Status = "closed"
	default:
		t.ts.Status = "idle"
	}

	if err != nil {
		t.updateStates(nil)
	}

	return err
}

func (t *table) startGame() error {

	// Preparing options
	var opts *pokerface.GameOptions

	// Preparing deck
	switch t.options.GameType {
	case "short_deck":
		opts = pokerface.NewStardardGameOptions()
		opts.Deck = pokerface.NewShortDeckCards()
	default:
		opts = pokerface.NewStardardGameOptions()
		opts.Deck = pokerface.NewStandardDeckCards()
	}

	// Preparing options
	opts.Ante = t.options.Ante
	opts.Blind.Dealer = t.options.Blind.Dealer
	opts.Blind.SB = t.options.Blind.SB
	opts.Blind.BB = t.options.Blind.BB

	// Preparing players
	seats := t.sm.GetPlayableSeats()
	for i, s := range seats {
		s.Player.GameIdx = i
		opts.Players = append(opts.Players, &pokerface.PlayerSetting{
			Bankroll:  s.Player.Bankroll,
			Positions: s.Player.Positions,
		})
	}

	// Create a new game with backend
	t.g = NewGame(t.b, opts)
	t.g.OnStateUpdated(func(gs *pokerface.GameState) {
		//fmt.Println(gs.Status.CurrentEvent.Name)
		t.updateStates(gs)
	})

	err := t.g.Start()
	if err != nil {
		return err
	}

	t.gameCount++
	t.ts.Status = "playing"

	return nil
}
