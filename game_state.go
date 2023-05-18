package main

import "github.com/cfsghost/pokerface/combination"

type GameState struct {
	GameID    string        `json:"game_id"`
	CreatedAt int64         `json:"created_at"`
	UpdatedAt int64         `json:"updated_at"`
	Meta      Meta          `json:"meta"`
	Status    Status        `json:"status"`
	Players   []PlayerState `json:"players"`
	Result    Result        `json:"result"`
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

type WorkflowEvent struct {
	Name    string      `json:"name"`
	Runtime interface{} `json:"runtime"`
}

type Status struct {
	MiniBet             int64          `json:"min_bet"`
	Pots                map[string]Pot `json:"pots"`
	Round               string         `json:"round"`
	Burned              []string       `json:"burned"`
	Board               []string       `json:"board"`
	CurrentDeckPosition int            `json:"current_deck_position"`
	CurrentWager        int64          `json:"current_wager"`
	CurrentRaiser       int            `json:"current_raiser"`
	CurrentPlayer       int            `json:"current_player"`
	CurrentEvent        *WorkflowEvent `json:"current_event"`
}

type Pot struct {
	Chips        int   `json:"chips"`
	Contributers []int `json:"contributers"`
}

type PlayerState struct {
	Idx              int             `json:"idx"`
	Position         []string        `json:"position"`
	DidAction        string          `json:"did_action"`
	Bankroll         int64           `json:"bankroll"`
	InitialStackSize int64           `json:"initial_stack_size"` // bankroll - pot
	StackSize        int64           `json:"stack_size"`         // initial_stack_size - wager
	Pot              int64           `json:"pot"`
	Wager            int64           `json:"wager"`
	HoleCards        []string        `json:"hole_cards"`
	Fold             bool            `json:"fold"`
	ActionCount      int             `json:"action_count"`
	Combination      CombinationInfo `json:"combination"`

	// Actions
	AllowedActions []string `json:"allowed_actions"`
}

type CombinationInfo struct {
	Type  string   `json:"type"`
	Cards []string `json:"cards"`
	Power int      `json:"power"`
}

type Result struct {
	Players []PlayerResult `json:"players"`
	Pots    []PotResult    `json:"pots"`
}

type PlayerResult struct {
	Idx     int   `json:"idx"`
	Final   int64 `json:"final"`
	Changed int64 `json:"changed"`
}

type PotResult struct {
	Chips   int64    `json:"chips"`
	Winners []Winner `json:"winners"`
}

type Winner struct {
	Idx   int   `json:"idx"`
	Chips int64 `json:"chips"`
}
