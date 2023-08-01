package competition

import (
	"errors"
	"sync"

	"github.com/weedbox/pokerface/match"
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
	TableManager() TableManager
	TableBackend() TableBackend
	Start() error
	Close() error
	GetOptions() *Options
	GetTableCount() int64
	Match() match.Match
	ReserveSeat(tableID string, seatID int, playerID string) (int, error)
}

type competition struct {
	options    *Options
	tm         TableManager
	tb         TableBackend
	m          match.Match
	mtb        match.TableBackend
	players    []*PlayerInfo
	s          *State
	isRunning  bool
	isJoinable bool
	mu         sync.RWMutex
}

type CompetitionOpt func(*competition)

func WithTableBackend(tb TableBackend) CompetitionOpt {
	return func(c *competition) {
		c.tb = tb
	}
}

func WithMatchBackend(m match.Match) CompetitionOpt {
	return func(c *competition) {
		c.m = m
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

	// Table backend
	if c.tb == nil {
		c.tb = NewNativeTableBackend(table.NewNativeBackend())
	}

	// Table backend of match
	if c.mtb == nil {
		c.mtb = NewNativeMatchTableBackend(c)
	}

	if c.m == nil {

		// Initializing match
		opts := match.NewOptions()
		opts.WaitingPeriod = c.GetOptions().TableAllocationPeriod
		opts.MaxTables = c.GetOptions().MaxTables
		opts.MaxSeats = c.GetOptions().Table.MaxSeats

		c.m = match.NewMatch(
			opts,
			match.WithTableBackend(NewNativeMatchTableBackend(c)),
		)
	}

	c.tm = NewTableManager(options, c.tb)

	return c
}

func (c *competition) getPlayerByID(playerID string) (*PlayerInfo, error) {

	for _, p := range c.players {
		if p.ID == playerID {
			return p, nil
		}
	}

	return nil, ErrNotFoundPlayer
}

func (c *competition) TableManager() TableManager {
	return c.tm
}

func (c *competition) TableBackend() TableBackend {
	return c.tb
}

func (c *competition) Match() match.Match {
	return c.m
}

func (c *competition) GetOptions() *Options {
	return c.options
}

func (c *competition) GetTableCount() int64 {
	return c.tm.GetTableCount()
}

func (c *competition) GetPlayers() []*PlayerInfo {
	return c.players
}

func (c *competition) GetPlayerByID(playerID string) (*PlayerInfo, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.getPlayerByID(playerID)
}

func (c *competition) GetPlayerIndexByID(playerID string) (int, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	for i, p := range c.players {
		if p.ID == playerID {
			return i, nil
		}
	}

	return -1, ErrNotFoundPlayer
}

func (c *competition) Register(playerID string, bankroll int64) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isJoinable {
		return ErrNotJoinable
	}

	_, err := c.getPlayerByID(playerID)
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
		err := c.m.Join(playerID)
		if err != nil {
			return err
		}

		p.Participated = true
	}

	return nil
}

func (c *competition) Unregister(playerID string) error {

	c.mu.Lock()
	defer c.mu.Unlock()

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

		err := c.m.Join(p.ID)
		if err != nil {
			return err
		}

		p.Participated = true
	}

	return nil
}

func (c *competition) Close() error {
	return c.m.Close()
}

func (c *competition) BuyIn(p *PlayerInfo) error {

	_, err := c.GetPlayerByID(p.ID)
	if err != ErrNotFoundPlayer {
		// Existing already
		return err
	}

	// Allocate seat
	err = c.m.Join(p.ID)
	if err != nil {
		return err
	}

	p.Participated = true

	return nil
}

func (c *competition) ReserveSeat(tableID string, seatID int, playerID string) (int, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	p, err := c.getPlayerByID(playerID)
	if err != nil {
		return -1, err
	}

	return c.tm.ReserveSeat(tableID, seatID, p)
}
