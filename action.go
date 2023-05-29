package pokerface

func (g *game) ReadyForAll() error {

	for _, p := range g.GetPlayers() {
		err := p.Ready()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *game) Ready(playerIdx int) error {
	return g.Player(playerIdx).Ready()
}

func (g *game) PayAnte() error {

	if g.gs.Meta.Ante == 0 {
		return nil
	}

	for _, p := range g.GetPlayers() {
		err := p.Pay(g.gs.Meta.Ante)
		if err != nil {
			return err
		}
	}

	return nil
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
