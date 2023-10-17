package pokerface

import (
	"fmt"

	"github.com/weedbox/pokerface/pot"
)

func (g *game) updatePots() error {

	ll := pot.NewLevelList()

	for _, p := range g.gs.Players {
		ll.AddContributor(p.Pot+p.Wager, p.Idx, p.Fold)
	}

	g.gs.Status.Pots = ll.GetPots()

	return nil
}

func (g *game) PrintPots() {

	for _, p := range g.GetState().Status.Pots {
		fmt.Println("======= POT")
		fmt.Println("Contributors", p.Contributors)
		fmt.Println("Level", p.Level)
		fmt.Println("Wager", p.Wager)
		fmt.Println("Total", p.Total)
	}
}
