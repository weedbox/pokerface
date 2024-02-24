package pokerface

import (
	"github.com/weedbox/pokerface/combination"
	"github.com/weedbox/pokerface/pot"
	"github.com/weedbox/pokerface/settlement"
)

type GameState struct {
	GameID    string             `json:"game_id"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Meta      Meta               `json:"meta"`
	Status    Status             `json:"status"`
	Players   []*PlayerState     `json:"players"`
	Result    *settlement.Result `json:"result,omitempty"`
}

type Meta struct {
	Ante                   int64                     `json:"ante"`
	Blind                  BlindSetting              `json:"blind"`
	Limit                  string                    `json:"limit"`
	HoleCardsCount         int                       `json:"hole_cards_count"`
	RequiredHoleCardsCount int                       `json:"required_hole_cards_count"`
	CombinationPowers      combination.PowerRankings `json:"combination_powers"`
	Deck                   []string                  `json:"deck"`
	BurnCount              int                       `json:"burn_count"`
}

type Action struct {
	Source int    `json:"source"`
	Type   string `json:"type"`
	Value  int64  `json:"value,omitempty"`
}

type Status struct {
	MiniBet             int64      `json:"mini_bet"`
	Pots                []*pot.Pot `json:"pots"`
	Round               string     `json:"round,omitempty"`
	Burned              []string   `json:"burned,omitempty"`
	Board               []string   `json:"board,omitempty"`
	PreviousRaiseSize   int64      `json:"previous_raise_size"`
	CurrentDeckPosition int        `json:"current_deck_position"`
	CurrentRoundPot     int64      `json:"current_round_pot"`
	CurrentWager        int64      `json:"current_wager"`
	CurrentRaiser       int        `json:"current_raiser"`
	CurrentPlayer       int        `json:"current_player"`
	CurrentEvent        string     `json:"current_event"`
	LastAction          *Action    `json:"last_action,omitempty"`
}

type PlayerState struct {
	Idx       int      `json:"idx"`
	Positions []string `json:"positions"`

	// Status
	Acted          bool     `json:"acted"`
	DidAction      string   `json:"did_action,omitempty"`
	Fold           bool     `json:"fold"`
	VPIP           bool     `json:"vpip"` // Voluntarily Put In Pot
	AllowedActions []string `json:"allowed_actions,omitempty"`

	// Stack and wager
	Bankroll         int64 `json:"bankroll"`
	InitialStackSize int64 `json:"initial_stack_size"` // bankroll - pot
	StackSize        int64 `json:"stack_size"`         // initial_stack_size - wager
	Pot              int64 `json:"pot"`
	Wager            int64 `json:"wager"`

	// Hole cards information
	HoleCards   []string         `json:"hole_cards,omitempty"`
	Combination *CombinationInfo `json:"combination,omitempty"`
}

type CombinationInfo struct {
	Type  string   `json:"type"`
	Cards []string `json:"cards"`
	Power int      `json:"power"`
}

func (gs *GameState) AsPlayer(idx int) {

	gs.Meta.Deck = []string{}
	gs.Status.Burned = []string{}

	// Do nothing if game has been closed already
	if gs.Status.CurrentEvent == "GameClosed" {

		for _, p := range gs.Players {
			if p.Idx == idx {
				continue
			}

			// Hide private information if player do fold
			if p.Fold {
				p.HoleCards = []string{}
				p.Combination = nil
			}
		}

		return
	}

	for _, p := range gs.Players {
		if p.Idx == idx {
			continue
		}

		// Hide private information
		p.HoleCards = []string{}
		p.Combination = nil
	}
}

func (gs *GameState) AsObserver() {

	gs.Meta.Deck = []string{}
	gs.Status.Burned = []string{}

	if gs.Status.CurrentEvent == "GameClosed" {

		for _, p := range gs.Players {

			// Hide private information if player do fold
			if p.Fold {
				p.HoleCards = []string{}
				p.Combination = nil
			}
		}

		return
	}

	// Hide all private information
	for _, p := range gs.Players {
		p.HoleCards = []string{}
		p.Combination = nil
	}
}

func (gs *GameState) GetPlayer(idx int) *PlayerState {

	if idx < 0 || idx >= len(gs.Players) {
		return nil
	}

	return gs.Players[idx]
}

func (gs *GameState) HasPosition(idx int, position string) bool {

	p := gs.GetPlayer(idx)
	if p == nil {
		return false
	}

	for _, pos := range p.Positions {
		if pos == position {
			return true
		}
	}

	return false
}

func (gs *GameState) HasAction(idx int, action string) bool {

	p := gs.GetPlayer(idx)
	if p == nil {
		return false
	}

	for _, aa := range p.AllowedActions {
		if aa == action {
			return true
		}
	}

	return false
}

func (ps *PlayerState) AllowAction(action string) {

	for _, aa := range ps.AllowedActions {
		if aa == action {
			return
		}
	}

	ps.AllowedActions = append(ps.AllowedActions, action)
}
