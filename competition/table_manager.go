package competition

import (
	"sync/atomic"

	"github.com/weedbox/pokerface/table"
)

type TableManager interface {
	Initialize() error
	CreateTable() (*table.State, error)
	ReleaseTable(tableID string) error
	ActivateTable(tableID string) error
	GetTables() []*table.State
	GetTableState(tableID string) *table.State
	GetTableCount() int64
	UpdateTableState(ts *table.State) error
	ReserveSeat(tableID string, seatID int, p *PlayerInfo) (int, error)
	OnTableStateUpdated(fn func(ts *table.State))
}

type tableManager struct {
	options             *Options
	b                   TableBackend
	tables              map[string]*table.State
	count               int64
	onTableStateUpdated func(ts *table.State)
}

func NewTableManager(options *Options, b TableBackend) TableManager {
	return &tableManager{
		options:             options,
		b:                   b,
		tables:              make(map[string]*table.State),
		onTableStateUpdated: func(ts *table.State) {},
	}
}

func (tm *tableManager) Initialize() error {

	return nil
	/*
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
	*/
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

	atomic.AddInt64(&tm.count, 1)

	return ts, nil
}

func (tm *tableManager) ReleaseTable(tableID string) error {

	err := tm.b.ReleaseTable(tableID)
	if err != nil {
		return err
	}

	atomic.AddInt64(&tm.count, -1)

	return nil
}

func (tm *tableManager) ActivateTable(tableID string) error {
	return tm.b.ActivateTable(tableID)
}

func (tm *tableManager) UpdateTableState(ts *table.State) error {

	_, ok := tm.tables[ts.ID]
	if !ok {
		return ErrNotFoundTable
	}

	tm.tables[ts.ID] = ts

	tm.onTableStateUpdated(ts)

	return nil
}

func (tm *tableManager) GetTableState(tableID string) *table.State {

	ts, ok := tm.tables[tableID]
	if !ok {
		return nil
	}

	return ts
}

func (tm *tableManager) GetTables() []*table.State {

	tables := make([]*table.State, 0)

	for _, ts := range tm.tables {
		tables = append(tables, ts)
	}

	return tables
}

func (tm *tableManager) GetTableCount() int64 {
	return atomic.LoadInt64(&tm.count)
}

func (tm *tableManager) ReserveSeat(tableID string, seatID int, p *PlayerInfo) (int, error) {
	return tm.b.ReserveSeat(tableID, -1, p)
}

func (tm *tableManager) OnTableStateUpdated(fn func(ts *table.State)) {
	tm.onTableStateUpdated = fn
}
