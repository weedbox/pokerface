package competition

import (
	"github.com/weedbox/pokerface/table"
)

type TableManager interface {
	Initialize() error
	CreateTable() (*table.State, error)
	ActivateTable(tableID string) error
	GetTableState(tableID string) *table.State
	GetTableCount() int
	DispatchPlayer(p *PlayerInfo) error
}

type tableManager struct {
	options *Options
	b       TableBackend
	tables  map[string]*table.State
}

func NewTableManager(options *Options, b TableBackend) *tableManager {
	return &tableManager{
		options: options,
		b:       b,
		tables:  make(map[string]*table.State),
	}
}

func (tm *tableManager) Initialize() error {

	// Only one static table
	if tm.options.MaxTables == 1 {
		ts, err := tm.CreateTable()
		if err != nil {
			return err
		}

		return tm.ActivateTable(ts.ID)
	}

	//TODO: Initializing dynamic allocation

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
	for _, ts := range tm.tables {
		p.Participated = true
		return tm.b.ReserveSeat(ts.ID, -1, p)
	}

	return nil
}
