package main

import "github.com/cfsghost/pokerface/waitgroup"

func (g *game) WaitForAllPlayersReady() (*waitgroup.WaitGroup, error) {

	// First time to initializing waitgroup
	if g.gs.Status.CurrentEvent.Runtime == nil {

		// Preparing runtime to wait for player ready
		players := make([]int, 0, len(g.gs.Players))
		for _, p := range g.gs.Players {
			players = append(players, p.Idx)
		}

		r := waitgroup.NewWaitReadyRuntime(players)

		g.gs.Status.CurrentEvent.Runtime = r
	}

	// Getting runtime from state
	r := g.gs.Status.CurrentEvent.Runtime.(*waitgroup.WaitGroupRuntime)

	// Initializing wait group based on current event runtime
	wg := waitgroup.NewWaitGroup(waitgroup.TypeReady, r, waitgroup.WaitReady)

	// Update states based on wait group
	for _, p := range g.gs.Players {

		s := wg.GetStateByIdx(p.Idx)
		if s == nil {
			continue
		}

		if s.State {
			continue
		}

		// this player doesn't get ready yet
		g.Player(p.Idx).AllowActions([]string{
			"ready",
		})
	}

	g.wg = wg

	return wg, nil
}

func (g *game) WaitForAllPlayersPaidAnte() (*waitgroup.WaitGroup, error) {

	// First time to initializing waitgroup
	if g.gs.Status.CurrentEvent.Runtime == nil {

		// Preparing runtime to wait for ante
		players := make([]int, 0, len(g.gs.Players))
		for _, p := range g.gs.Players {
			players = append(players, p.Idx)
		}

		r := waitgroup.NewWaitPayAnteRuntime(players)

		g.gs.Status.CurrentEvent.Runtime = r
	}

	// Getting runtime from state
	r := g.gs.Status.CurrentEvent.Runtime.(*waitgroup.WaitGroupRuntime)

	// Initializing wait group based on current event runtime
	wg := waitgroup.NewWaitGroup(waitgroup.TypePayAnte, r, waitgroup.WaitReady)

	// Update states based on wait group
	for _, p := range g.gs.Players {

		s := wg.GetStateByIdx(p.Idx)
		if s == nil {
			continue
		}

		if s.State {
			continue
		}

		// this player doesn't get ready yet
		g.Player(p.Idx).AllowActions([]string{
			"pay_ante",
		})
	}

	g.wg = wg

	return wg, nil
}
