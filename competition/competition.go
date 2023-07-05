package competition

import (
	"errors"

	"github.com/weedbox/pokerface/table"
)

var (
	ErrRunningAlready      = errors.New("competition: running already")
	ErrNotJoinable         = errors.New("competition: not joinable")
	ErrParticipatedAlready = errors.New("competition: player participated already")
	ErrNotFoundTable       = errors.New("competition: not found table")
	ErrNotFoundPlayer      = errors.New("competition: not found player")
	ErrPlayerExistsAlready = errors.New("competition: player exists already")
)

type Competition interface {
}

type competition struct {
	options    *Options
	tm         TableManager
	tb         TableBackend
	players    []*PlayerInfo
	s          *State
	isRunning  bool
	isJoinable bool
}

type CompetitionOpt func(*competition)

func WithTableBackend(tb TableBackend) CompetitionOpt {
	return func(c *competition) {
		c.tb = tb
	}
}

func NewCompetition(options *Options, opts ...CompetitionOpt) *competition {

	c := &competition{
		options:    options,
		players:    make([]*PlayerInfo, 0),
		s:          NewState(),
		isJoinable: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.tb == nil {
		c.tb = NewNativeTableBackend(table.NewNativeBackend())
	}

	c.tm = NewTableManager(options, c.tb)

	return c
}

func (c *competition) GetTableCount() int {
	return c.tm.GetTableCount()
}

func (c *competition) GetPlayers() []*PlayerInfo {
	return c.players
}

func (c *competition) GetPlayerByID(playerID string) (*PlayerInfo, error) {

	for _, p := range c.players {
		if p.ID == playerID {
			return p, nil
		}
	}

	return nil, ErrNotFoundPlayer
}

func (c *competition) GetPlayerIndexByID(playerID string) (int, error) {

	for i, p := range c.players {
		if p.ID == playerID {
			return i, nil
		}
	}

	return -1, ErrNotFoundPlayer
}

func (c *competition) Register(playerID string, bankroll int64) error {

	if !c.isJoinable {
		return ErrNotJoinable
	}

	_, err := c.GetPlayerByID(playerID)
	if err != ErrNotFoundPlayer {
		// Existing already
		return ErrPlayerExistsAlready
	}

	p := &PlayerInfo{
		ID:       playerID,
		Bankroll: bankroll,
	}

	c.players = append(c.players, p)

	// Dispatch player if competition is running already
	if c.isRunning {
		c.tm.DispatchPlayer(p)
	}

	return nil
}

func (c *competition) Unregister(playerID string) error {

	idx := -1
	var found *PlayerInfo
	for i, p := range c.players {
		if p.ID == playerID {
			idx = i
			found = p
			break
		}
	}

	if found == nil {
		return ErrNotFoundPlayer
	}

	// Disallow to unregister if player participated game
	if found.Participated {
		return ErrParticipatedAlready
	}

	// Remove player from list
	c.players = append(c.players[:idx], c.players[idx+1:]...)

	return nil
}

func (c *competition) Start() error {

	if c.isRunning {
		return ErrRunningAlready
	}

	// Initializing tables
	err := c.tm.Initialize()
	if err != nil {
		return err
	}

	c.isRunning = true

	// Dispatching registered players who is waiting for game start
	for _, p := range c.players {

		// Participated already
		if p.Participated {
			continue
		}

		c.tm.DispatchPlayer(p)
	}

	return nil
}

func (c *competition) BuyIn(p *PlayerInfo) error {

	_, err := c.GetPlayerByID(p.ID)
	if err != ErrNotFoundPlayer {
		// Existing already
		return err
	}

	// Allocate seat
	return c.tm.DispatchPlayer(p)
}
