package competition

import (
	"github.com/weedbox/pokerface/table"
)

type TableInfo struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type TableBackend interface {
	CreateTable(opts *table.Options) (*table.State, error)
	ActivateTable(tableID string) error
	SetJoinable(tableID string, isJoinable bool) error
	ReleaseTable(tableID string) error
	ReserveSeat(tableID string, seatID int, player *PlayerInfo) (int, error)
	OnTableUpdated(fn func(ts *table.State))
}
