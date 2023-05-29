package pokerface

import (
	"github.com/cfsghost/pokerface/combination"
	"github.com/cfsghost/pokerface/pot"
	"github.com/cfsghost/pokerface/settlement"
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

type Event struct {
	Name    string        `json:"name,omitempty"`
	Payload *EventPayload `json:"payload,omitempty"`
}

type Status struct {
	MiniBet             int64      `json:"min_bet"`
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
	CurrentEvent        *Event     `json:"current_event"`
}

type PlayerState struct {
	Idx              int             `json:"idx"`
	Positions        []string        `json:"positions"`
	DidAction        string          `json:"did_action,omitempty"`
	Bankroll         int64           `json:"bankroll"`
	InitialStackSize int64           `json:"initial_stack_size"` // bankroll - pot
	StackSize        int64           `json:"stack_size"`         // initial_stack_size - wager
	Pot              int64           `json:"pot"`
	Wager            int64           `json:"wager"`
	HoleCards        []string        `json:"hole_cards,omitempty"`
	Fold             bool            `json:"fold"`
	ActionCount      int             `json:"action_count"`
	Combination      CombinationInfo `json:"combination,omitempty"`

	// Actions
	AllowedActions []string `json:"allowed_actions,omitempty"`
}

type CombinationInfo struct {
	Type  string   `json:"type"`
	Cards []string `json:"cards"`
	Power int      `json:"power"`
}
