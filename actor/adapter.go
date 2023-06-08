package actor

import pokertable "github.com/weedbox/pokertable/model"

type Adapter interface {
	SetActor(a Actor)
	UpdateTableState(t *pokertable.Table) error

	// Player actions
	Ready(playerID string) error
	Pay(playerID string, chips int64) error
	Check(playerID string) error
	Bet(playerID string, chips int64) error
	Call(playerID string) error
	Fold(playerID string) error
	Allin(playerID string) error
	Raise(playerID string, chipLevel int64) error
}
