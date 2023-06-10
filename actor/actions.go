package actor

type Actions interface {
	Pass() error
	Ready() error
	Pay(chips int64) error
	Check() error
	Bet(chips int64) error
	Call() error
	Fold() error
	Allin() error
	Raise(chipLevel int64) error
}

type actions struct {
	actor    Actor
	playerID string
}

func NewActions(actor Actor, playerID string) Actions {
	return &actions{
		actor:    actor,
		playerID: playerID,
	}
}

func (a *actions) Pass() error {
	return a.actor.GetTable().Pass(a.playerID)
}

func (a *actions) Ready() error {
	return a.actor.GetTable().Ready(a.playerID)
}

func (a *actions) Pay(chips int64) error {
	return a.actor.GetTable().Pay(a.playerID, chips)
}

func (a *actions) Check() error {
	return a.actor.GetTable().Check(a.playerID)
}

func (a *actions) Bet(chips int64) error {
	return a.actor.GetTable().Bet(a.playerID, chips)
}

func (a *actions) Call() error {
	return a.actor.GetTable().Call(a.playerID)
}

func (a *actions) Fold() error {
	return a.actor.GetTable().Fold(a.playerID)
}

func (a *actions) Allin() error {
	return a.actor.GetTable().Allin(a.playerID)
}

func (a *actions) Raise(chipLevel int64) error {
	return a.actor.GetTable().Raise(a.playerID, chipLevel)
}
