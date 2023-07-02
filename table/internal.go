package table

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/weedbox/pokerface"
)

func (t *table) tableLoop() {

	for interval := range t.gameLoop {

		err := t.delay(interval, func() error {
			return t.nextGame()
		})

		switch err {
		case ErrMaxGamesExceeded:
			t.ts.Status = "closed"
			t.updateStates(nil)
			return
		case ErrTimesUp:
			t.ts.Status = "closed"
			t.updateStates(nil)
			return
		}

		if err := t.checkEndConditions(); err != nil {
			t.ts.Status = "closed"
			t.updateStates(nil)
			return
		}

		switch err {
		case ErrInsufficientNumberOfPlayers:
			// Waiting for more players
			t.ts.Status = "idle"
			continue
		case ErrGameCancelled:
			// Do nothing
			continue
		}

		// Continue to the next game
		t.ts.Status = "pending"
		t.NewGame(t.options.Interval)
	}

	t.isRunning = false
}

func (t *table) delay(interval int, fn func() error) error {

	var err error
	var wg sync.WaitGroup
	wg.Add(1)

	t.tb.NewTask(time.Duration(interval)*time.Second, func(isCancelled bool) {

		defer wg.Done()

		if isCancelled {
			err = ErrGameCancelled
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

	// Event
	t.onStateUpdated(&state)

	return nil
}

func (t *table) checkEndConditions() error {

	if t.options.MaxGames > 0 && t.options.MaxGames == t.gameCount {
		return ErrMaxGamesExceeded
	}

	// Check remaining time
	if time.Now().Unix() >= t.ts.EndTime {
		// Times up!
		return ErrTimesUp
	}

	return nil
}

func (t *table) nextGame() error {

	if t.isPaused {
		return ErrGameCancelled
	}

	t.ts.Status = "preparing"

	if err := t.checkEndConditions(); err != nil {
		return err
	}

	err := t.setupPosition()
	if err != nil {
		return err
	}

	// Check the number of player
	playableCount := t.sm.GetPlayableSeatCount()
	if t.gameCount == 0 && playableCount < t.options.InitialPlayers {
		return ErrInsufficientNumberOfPlayers
	} else if playableCount < t.options.MinPlayers {
		return ErrInsufficientNumberOfPlayers
	}

	return t.startGame()
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
		//fmt.Println(gs.Status.CurrentEvent)
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
