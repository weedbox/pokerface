package pokerface

func (g *game) ReadyForAll() error {

	if g.gs.Status.CurrentEvent != "ReadyRequested" {
		return ErrInvalidAction
	}

	return g.EmitEvent(GameEvent_Readiness)
}

func (g *game) Pass() error {
	return g.GetCurrentPlayer().Pass()
}

func (g *game) PayAnte() error {

	if g.gs.Meta.Ante == 0 {
		return ErrInvalidAction
	}

	if g.gs.Status.CurrentEvent != "AnteRequested" {
		return ErrInvalidAction
	}

	for _, p := range g.GetPlayers() {
		err := p.PayAnte()
		if err != nil {
			return err
		}
	}

	return g.EmitEvent(GameEvent_AntePaid)
}

func (g *game) PayBlinds() error {

	if g.gs.Status.CurrentEvent != "BlindsRequested" {
		return ErrInvalidAction
	}

	for _, p := range g.GetPlayers() {
		err := p.PayBlinds()
		if err != nil {
			return err
		}
	}

	// Minimal raise size
	if g.gs.Meta.Blind.BB > 0 {
		g.gs.Status.PreviousRaiseSize = g.gs.Meta.Blind.BB
	} else {
		g.gs.Status.PreviousRaiseSize = g.gs.Meta.Blind.Dealer
	}

	return g.EmitEvent(GameEvent_BlindsPaid)
}

func (g *game) Pay(chips int64) error {
	return g.GetCurrentPlayer().Pay(chips)
}

func (g *game) Fold() error {
	return g.GetCurrentPlayer().Fold()
}

func (g *game) Check() error {
	return g.GetCurrentPlayer().Check()
}

func (g *game) Call() error {
	return g.GetCurrentPlayer().Call()
}

func (g *game) Allin() error {
	return g.GetCurrentPlayer().Allin()
}

func (g *game) Bet(chips int64) error {
	return g.GetCurrentPlayer().Bet(chips)
}

func (g *game) Raise(chipLevel int64) error {
	return g.GetCurrentPlayer().Raise(chipLevel)
}
