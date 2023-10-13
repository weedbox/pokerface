package pokerface

import (
	"fmt"

	"github.com/weedbox/pokerface/pot"
)

func (g *game) updatePots() error {

	pots := pot.NewPotList()

	for _, p := range g.gs.Players {

		// Not contributer
		if p.Wager == 0 {
			continue
		}

		pots.AddContributor(p.Wager, p.Idx)
	}

	// Merge pots into original pot list
	for i, pot := range pots.GetPots() {

		// More side pots or no main pot
		if i > 0 || len(g.gs.Status.Pots) == 0 {
			g.gs.Status.Pots = append(g.gs.Status.Pots, pot)
			continue
		}

		// Getting the last pot
		lastPot := g.gs.Status.Pots[len(g.gs.Status.Pots)-1]

		// Check contributors
		sameContributors := true
		for _, cIdx := range lastPot.Contributors {

			// the lists of contributors are different
			if !pot.ContributorExists(cIdx) {
				sameContributors = false
				break
			}
		}

		if !sameContributors {
			// Do not merge
			g.gs.Status.Pots = append(g.gs.Status.Pots, pot)
			continue
		}

		// Merge pot
		lastPot.Level += pot.Level
		lastPot.Wager += pot.Wager
		lastPot.Total += pot.Total
	}

	// After considering the folded players, arrange all the pots
	for _, pot := range g.gs.Status.Pots {

		contributors := make([]int, 0)

		for _, cIdx := range pot.Contributors {

			// Remove folded player
			if g.Player(cIdx).State().Fold {
				continue
			}

			contributors = append(contributors, cIdx)
		}

		pot.Contributors = contributors
	}

	// Merge pots which contains the same contritutors
	oldPots := g.gs.Status.Pots
	g.gs.Status.Pots = make([]*pot.Pot, 0)
	var prevPot *pot.Pot = nil
	for i, pot := range oldPots {

		if i == 0 {
			g.gs.Status.Pots = append(g.gs.Status.Pots, pot)
			prevPot = pot
			continue
		}

		if len(prevPot.Contributors) != len(pot.Contributors) {
			g.gs.Status.Pots = append(g.gs.Status.Pots, pot)
			prevPot = pot
			continue
		}

		prevPot.Level += pot.Level
		prevPot.Wager += pot.Wager
		prevPot.Total += pot.Total
	}

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
