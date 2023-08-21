package table

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/weedbox/pokerface"
	"github.com/weedbox/pokerface/seat_manager"
)

func (t *table) tableLoop() {

	//count := 0
	for interval := range t.gameLoop {

		if !t.isRunning {
			fmt.Println("TABLE NotRunning", t.ts.ID)
			break
		}

		//count++
		//fmt.Println("tableLoop", t.GetState().ID, count)

		err := t.delay(interval, func() error {
			return t.prepareNextGame()
		})

		switch err {
		case ErrMaxGamesExceeded:
			fmt.Println("TABLE ErrMaxGamesExceeded")
			t.ts.Status = "closed"
			t.updateGameState(nil)
			return
		case ErrTimesUp:
			fmt.Println("TABLE ErrTimesUp")
			t.ts.Status = "closed"
			t.updateGameState(nil)
			return
		case ErrInsufficientNumberOfPlayers:

			// Nobody can join so so table should be closed
			if !t.options.Joinable {
				fmt.Println("TABLE ErrInsufficientNumberOfPlayers")
				t.ts.Status = "closed"
				t.updateGameState(nil)
				return
			}

			// Waiting for more players
			t.ts.Status = "idle"
			continue
		case ErrGameCancelled:
			// Do nothing
			continue
		}

		if !t.isRunning {
			fmt.Println("TABLE NotRunning", t.ts.ID)
			break
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

	t.tb.NewTask(time.Duration(interval)*time.Millisecond, func(isCancelled bool) {

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

	// No need to prepare game roles as the roles are already in position
	if t.inPosition {
		return nil
	}

	// Calculating positions for players
	err := t.sm.Next()
	if err != nil {

		if err == seat_manager.ErrInsufficientNumberOfPlayers {
			return ErrInsufficientNumberOfPlayers
		}

		return err
	}

	// Updating seat and position information for players
	t.mu.RLock()
	t.ts.ResetPositions()
	t.mu.RUnlock()

	seats := t.sm.GetSeats()
	for _, s := range seats {

		if s.Player == nil {
			continue
		}

		p := t.GetPlayerByID(s.Player.(*PlayerInfo).ID)
		p.Playable = false

		// Update position
		positions := make([]string, 0)

		if s == t.sm.Dealer() {
			positions = append(positions, "dealer")
		}

		if s == t.sm.SmallBlind() {
			positions = append(positions, "sb")
		} else if s == t.sm.BigBlind() {
			positions = append(positions, "bb")
		}

		p.Positions = positions

		// Update states
		if !s.IsReserved && s.IsActive {
			p.Playable = true
		}
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

		p := ts.GetPlayerByGameIdx(rs.Idx)
		if p == nil {
			continue
		}

		p.Bankroll = rs.Final

		// Not actively kicking players, waiting for requests to make players leave the table
		if p.Bankroll == 0 {
			t.sm.Reserve(p.SeatID)

			if t.ts.Options.EliminateMode == "leave" {
				//fmt.Println("updatePlayerStates", ts.ID, "LEAVE", p.SeatID, p.ID)
				t.leave(p.SeatID)
			}
		}
	}

	return nil
}

func (t *table) emitStateUpdated() {
	state := t.cloneState()
	t.onStateUpdated(state)
}

func (t *table) cloneState() *State {
	return t.ts.Clone()
	/*
	   // clone table state
	   data, err := json.Marshal(t.ts)

	   	if err != nil {
	   		return nil
	   	}

	   var state State
	   json.Unmarshal(data, &state)

	   return &state
	*/
}

func (t *table) updateGameState(gs *pokerface.GameState) error {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.ts.GameState = gs
	t.updatePlayerStates(t.ts)
	t.emitStateUpdated()

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

func (t *table) prepareNextGame() error {

	//t.mu.Lock()
	//defer t.mu.Unlock()

	if t.isPaused {
		return ErrGameCancelled
	}

	if err := t.checkEndConditions(); err != nil {
		return err
	}

	err := t.setupPosition()
	if err != nil {
		return err
	}

	t.ts.GameState = nil
	t.ts.Status = "preparing"

	// Check the number of player
	playableCount := t.sm.GetPlayableSeatCount()
	if t.gameCount == 0 && playableCount < t.options.InitialPlayers {
		return ErrInsufficientNumberOfPlayers
	} else if playableCount < t.options.MinPlayers {
		return ErrInsufficientNumberOfPlayers
	}

	err = t.startGame()
	if err != nil {
		return err
	}

	// Check conditions again
	if err := t.checkEndConditions(); err != nil {
		return err
	}

	// Preparing new positions
	err = t.setupPosition()
	if err != nil {
		return err
	}

	return nil
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

	// Clean legacy status
	t.mu.RLock()
	for _, p := range t.ts.Players {
		p.GameIdx = -1
	}
	t.mu.RUnlock()

	// Preparing players
	seats := t.sm.GetPlayableSeats()
	for i, s := range seats {
		s.Player.(*PlayerInfo).GameIdx = i
		opts.Players = append(opts.Players, &pokerface.PlayerSetting{
			Bankroll:  s.Player.(*PlayerInfo).Bankroll,
			Positions: s.Player.(*PlayerInfo).Positions,
		})
	}

	// Create a new game with backend
	t.g = NewGame(t.b, opts)

	// Preparing context
	ctx, cancel := context.WithCancel(context.Background())

	t.g.OnStateUpdated(func(gs *pokerface.GameState) {

		//fmt.Println(gs.GameID, gs.Status.CurrentEvent)
		t.updateGameState(gs)

		if gs.Status.CurrentEvent == "GameClosed" {
			cancel()
		}
	})

	err := t.g.Start()
	if err != nil {
		return err
	}

	t.gameCount++
	t.ts.Status = "playing"
	/*
		fmt.Println("startGame")
		for _, p := range t.ts.Players {
			fmt.Printf("[table] player=%s, gameIdx=%d, p.SeatID=%d, bankroll=%d, playable=%v\n", p.ID, p.GameIdx, p.SeatID, p.Bankroll, p.Playable)
		}
	*/
	// Waiting for game closed
	<-ctx.Done()

	t.inPosition = false

	return nil
}
