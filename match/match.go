package match

import (
	"errors"
	"fmt"
	"math"
	"sync/atomic"
)

type MatchStatus int

const (
	MatchStatus_Normal = iota
	MatchStatus_AfterRegDeadline
)

var (
	ErrNotFoundTable          = errors.New("match: not found table")
	ErrNotFoundAvailableTable = errors.New("match: not found available table")
	ErrAfterRegDealline       = errors.New("match: the final registration time has passed")
	ErrMaxPlayersReached      = errors.New("match: the current number of players has reached the maximum limit")
)

type Match interface {
	Options() *Options
	Dispatcher() Dispatcher
	WaitingRoom() WaitingRoom
	QueueManager() QueueManager
	TableBackend() TableBackend
	Runner() Runner
	TableMap() TableMap
	Close() error

	GetStatus() MatchStatus
	IsLastTableStage() bool
	Join(playerID string) error
	BreakTable(tableID string) error
	ApplySeatChanges(tableID string, sc *SeatChanges) error
	GetPlayerCount() int64
	PrintTables()
	AllocateTableWithPlayers(players []string) error
}

type MatchOpt func(*match)

type match struct {
	options     *Options
	status      MatchStatus
	qm          QueueManager
	wr          WaitingRoom
	tm          TableMap
	d           Dispatcher
	tb          TableBackend
	r           Runner
	playerCount int64

	onTableBroken  func(Match, *Table)
	onPlayerJoined func(Match, *Table, int, string)
}

func WithQueueManager(qm QueueManager) MatchOpt {
	return func(m *match) {
		m.qm = qm
	}
}

func WithTableBackend(tb TableBackend) MatchOpt {
	return func(m *match) {
		m.tb = tb
	}
}

func WithRunner(r Runner) MatchOpt {
	return func(m *match) {
		m.r = r
	}
}

func WithPlayerJoinedCallback(fn func(Match, *Table, int, string)) MatchOpt {
	return func(m *match) {
		m.onPlayerJoined = fn
	}
}

func WithTableBrokenCallback(fn func(Match, *Table)) MatchOpt {
	return func(m *match) {
		m.onTableBroken = fn
	}
}

func NewMatch(options *Options, opts ...MatchOpt) Match {

	m := &match{
		options:        options,
		status:         MatchStatus_Normal,
		onTableBroken:  func(Match, *Table) {},
		onPlayerJoined: func(Match, *Table, int, string) {},
	}

	for _, o := range opts {
		o(m)
	}

	m.wr = NewWaitingRoom(m)

	// Table map
	m.tm = NewTableMap(m)
	m.tm.OnTableBroken(func(table *Table) {
		m.onTableBroken(m, table)
	})

	m.tm.OnPlayerLeft(func(table *Table, seatID int, playerID string) {
		atomic.AddInt64(&m.playerCount, -int64(1))
	})

	m.tm.OnPlayerJoined(func(table *Table, seatID int, playerID string) {
		m.onPlayerJoined(m, table, seatID, playerID)
	})

	if m.qm == nil {
		m.qm = NewNativeQueueManager()
		err := m.qm.Connect()
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	// Table backend
	if m.tb == nil {
		m.tb = NewDummyTableBackend()
	}

	m.tb.OnTableUpdated(func(tableID string, sc *SeatChanges) {
		m.ApplySeatChanges(tableID, sc)
	})

	// Runner
	if m.r == nil {
		m.r = NewNativeRunner()
	}

	// Dispatcher
	m.d = NewDispatcher(m)
	m.d.OnFailure(func(err error, playerID string) {
		fmt.Println(err)
	})
	m.d.Start()

	return m
}

func (m *match) Options() *Options {
	return m.options
}

func (m *match) QueueManager() QueueManager {
	return m.qm
}

func (m *match) WaitingRoom() WaitingRoom {
	return m.wr
}

func (m *match) Runner() Runner {
	return m.r
}

func (m *match) TableBackend() TableBackend {
	return m.tb
}

func (m *match) TableMap() TableMap {
	return m.tm
}

func (m *match) Dispatcher() Dispatcher {
	return m.d
}

func (m *match) Close() error {

	err := m.d.Close()
	if err != nil {
		return err
	}

	err = m.qm.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *match) GetStatus() MatchStatus {
	return m.status
}

func (m *match) GetPlayerCount() int64 {
	return atomic.LoadInt64(&m.playerCount)
}

func (m *match) Join(playerID string) error {

	opts := m.Options()
	if opts.MaxTables > 0 {
		if opts.MaxTables*opts.MaxSeats <= int(m.GetPlayerCount()) {
			return ErrMaxPlayersReached
		}
	}

	if m.status == MatchStatus_AfterRegDeadline {
		return nil
	}

	err := m.d.Dispatch(playerID)
	if err != nil {
		return err
	}

	atomic.AddInt64(&m.playerCount, int64(1))

	return nil
}

func (m *match) BreakTable(tableID string) error {
	return m.tm.BreakTable(tableID)
}

func (m *match) ApplySeatChanges(tableID string, sc *SeatChanges) error {
	return m.tm.ApplySeatChanges(tableID, sc)
}

func (m *match) IsLastTableStage() bool {

	if m.status == MatchStatus_AfterRegDeadline {

		tableCount := m.TableMap().Count()

		// Final table
		if tableCount == 1 {
			return true
		}

		totalPlayers := m.GetPlayerCount()

		// Only one table is required for the rest of players
		proposedTableCount := math.Ceil(float64(totalPlayers) / float64(m.options.MaxSeats))
		if proposedTableCount == 1 {
			return true
		}
	}

	return false
}

func (m *match) PrintTables() {

	tables, _ := m.TableMap().GetTables()

	fmt.Printf("Current tables(%d):\n", len(tables))

	i := 0
	playerCount := 0
	for _, table := range tables {

		i++
		playerCount += table.GetPlayerCount()
		fmt.Printf("  table(%d) id=%s, player_count=%d, status=%d\n",
			i,
			table.ID(),
			table.GetPlayerCount(),
			table.GetStatus(),
		)
	}

	fmt.Printf("Total players = %d\n", m.GetPlayerCount())
	fmt.Printf("Waiting for dispatching = %d\n", m.Dispatcher().GetPendingCount())
}

func (m *match) AllocateTableWithPlayers(players []string) error {

	// Create a new table for players
	_, err := m.tm.CreateTable(players)
	if err != nil {
		return err
	}

	return nil
}
