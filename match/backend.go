package match

import "github.com/weedbox/pokerface/match/psae"

type Backend interface {
	AllocateTable() (string, error)
	BreakTable(tableID string) error
	Join(tableID string, players []*psae.Player) ([]int, error)
}
