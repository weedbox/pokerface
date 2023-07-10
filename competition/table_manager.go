package competition

import (
	"github.com/weedbox/pokerface/table"
)

type TableManager interface {
	Initialize() error
	CreateTable() (*table.State, error)
	BreakTable(tableID string) error
	ActivateTable(tableID string) error
	GetTableState(tableID string) *table.State
	GetTableCount() int
	DispatchPlayer(p *PlayerInfo) error
	ReserveSeat(tableID string, seatID int, p *PlayerInfo) (int, error)
}

type tableManager struct {
	options *Options
	b       TableBackend
	m       MatchBackend
	tables  map[string]*table.State
}

func NewTableManager(options *Options, b TableBackend, m MatchBackend) TableManager {
	return &tableManager{
		options: options,
		b:       b,
		m:       m,
		tables:  make(map[string]*table.State),
	}
}

func (tm *tableManager) Initialize() error {

	// Only one static table
	if tm.options.MaxTables == 1 {

		// Create table immediately
		ts, err := tm.CreateTable()
		if err != nil {
			return err
		}

		// Then activate
		return tm.ActivateTable(ts.ID)
	}

	return nil
}

func (tm *tableManager) CreateTable() (*table.State, error) {

	opts := table.NewOptions()
	opts.GameType = tm.options.GameType
	opts.MaxSeats = tm.options.Table.MaxSeats
	opts.Ante = tm.options.Table.Ante
	opts.Blind.Dealer = tm.options.Table.Blind.Dealer
	opts.Blind.SB = tm.options.Table.Blind.SB
	opts.Blind.BB = tm.options.Table.Blind.BB

	ts, err := tm.b.CreateTable(opts)
	if err != nil {
		return nil, err
	}

	tm.tables[ts.ID] = ts

	return ts, nil
}

func (tm *tableManager) BreakTable(tableID string) error {
	return tm.b.BreakTable(tableID)
}

func (tm *tableManager) ActivateTable(tableID string) error {
	return tm.b.ActivateTable(tableID)
}

func (tm *tableManager) GetTableState(tableID string) *table.State {

	ts, ok := tm.tables[tableID]
	if !ok {
		return nil
	}

	return ts
}

func (tm *tableManager) GetTableCount() int {
	return len(tm.tables)
}

func (tm *tableManager) DispatchPlayer(p *PlayerInfo) error {

	//TODO: replace this temporary solution
	// Find the on table for dispatching player into it
	/*
		for _, ts := range tm.tables {
			p.Participated = true
			_, err := tm.ReserveSeat(ts.ID, -1, p)
			if err != nil {
				return err
			}
		}
	*/

	p.Participated = true

	return tm.m.DispatchPlayer(p.ID)
}

func (tm *tableManager) ReserveSeat(tableID string, seatID int, p *PlayerInfo) (int, error) {
	return tm.b.ReserveSeat(tableID, -1, p)
}
