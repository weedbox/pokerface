package competition

import (
	"github.com/weedbox/pokerface/table"
)

type NativeTableBackend struct {
	geb    table.Backend
	tables map[string]table.Table
}

func NewNativeTableBackend(geb table.Backend) TableBackend {
	return &NativeTableBackend{
		geb:    geb,
		tables: make(map[string]table.Table),
	}
}

func (ntb *NativeTableBackend) CreateTable(opts *table.Options) (*table.State, error) {

	t := table.NewTable(opts, table.WithBackend(ntb.geb))

	ntb.tables[t.GetState().ID] = t

	return t.GetState(), nil
}

func (ntb *NativeTableBackend) ActivateTable(tableID string) error {

	t, ok := ntb.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	return t.Start()
}

func (ntb *NativeTableBackend) BreakTable(tableID string) error {

	_, ok := ntb.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	// TODO: implementation of breaking table

	return nil
}

func (ntb *NativeTableBackend) ReserveSeat(tableID string, seatID int, p *PlayerInfo) error {

	t, ok := ntb.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	_, err := t.Join(seatID, &table.PlayerInfo{
		ID:       p.ID,
		Bankroll: p.Bankroll,
	})

	return err
}
