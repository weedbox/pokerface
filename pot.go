package pokerface

import (
	"github.com/cfsghost/pokerface/pot"
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

	return nil
}
