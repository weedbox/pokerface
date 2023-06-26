package table

import (
	"encoding/json"
	"time"

	"github.com/weedbox/pokerface"
)

func (t *table) setupPosition() error {

	if t.inPosition {
		return nil
	}

	// Calculating positions for players
	err := t.sm.Next()
	if err != nil {
		return err
	}

	t.inPosition = true

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

	go t.onStateUpdated(&state)

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
	/*
		var wg sync.WaitGroup
		wg.Add(1)
		err = t.tb.NewTask(time.Duration(delay)*time.Second, func(isCancelled bool) {

			if isCancelled {
				return
			}

			wg.Done()
		})

		wg.Wait()
	*/
	if delay > 0 {
		<-time.After(time.Duration(delay) * time.Second)
	}

	// Starting a new game
	return t.startGame()
}

func (t *table) nextGame(delay int) error {

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
	opts := pokerface.NewStardardGameOptions()
	opts.Ante = t.options.Ante
	opts.Blind.Dealer = t.options.Blind.Dealer
	opts.Blind.SB = t.options.Blind.SB
	opts.Blind.BB = t.options.Blind.BB

	// Preparing deck
	switch t.options.GameType {
	case "short_deck":
		opts.Deck = pokerface.NewShortDeckCards()
	default:
		opts.Deck = pokerface.NewStandardDeckCards()
	}

	// Preparing players
	seats := t.sm.GetPlayableSeats()
	for _, s := range seats {
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
		t.nextGame(t.options.Interval)
	})

	err := t.g.Start()
	if err != nil {
		return err
	}

	t.gameCount++
	t.ts.Status = "playing"

	return nil
}
