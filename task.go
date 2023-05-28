package pokerface

import (
	"github.com/cfsghost/pokerface/task"
)

func (g *game) AssertReadyTask(taskName string) task.Task {

	// Getting current event
	event := g.GetEvent()

	t := event.Payload.Task.GetTask(taskName)

	// task doesn't exist
	if t == nil {

		// Create a new task to wait for ready
		wr := task.NewWaitReady(taskName)
		wr.PreparePlayerStates(len(g.gs.Players))

		// Update event payload
		wr.OnUpdated(func() {

			// Update allowed actions of players based on task state
			players := wr.GetPayload().(map[int]bool)
			for idx, isReady := range players {

				if isReady {
					g.Player(idx).AllowActions([]string{})
					continue
				}

				// Not ready so we are waiting for this player
				g.Player(idx).AllowActions([]string{
					"ready",
				})
			}
		})

		// Reset player action states when task completed
		wr.OnCompleted(func() {
			g.ResetAllPlayerAllowedActions()
		})

		event.Payload.Task.AddTask(wr)

		return wr
	}

	return t
}

func (g *game) WaitForAllPlayersReady(taskName string) bool {

	// Getting current event
	event := g.GetEvent()

	t := g.AssertReadyTask(taskName)

	// Getting current task to execute if it is what we need
	if event.Payload.Task.GetAvailableTask() != t {
		return false
	}

	if !t.Execute() {
		//fmt.Println("Waiting for ready")
		return false
	}

	return true
}

func (g *game) AssertPaymentTask(taskName string) task.Task {

	// Getting current event
	event := g.GetEvent()

	t := event.Payload.Task.GetTask(taskName)

	// task doesn't exist
	if t == nil {

		// Create a new task to wait for payment
		wp := task.NewWaitPay(taskName, g.gs.Meta.Ante)

		// All players must pay ante
		states := make([]int, 0)
		for _, ps := range g.gs.Players {
			states = append(states, ps.Idx)
		}
		wp.PrepareStates(states)

		// Update event payload
		wp.OnUpdated(func() {

			// Update allowed actions of players based on task state
			pr := wp.GetPayload().(task.PaymentRequest)
			for idx, chips := range pr.Players {

				// Paid already
				if chips > 0 {
					g.Player(idx).AllowActions([]string{})
					continue
				}

				// Not ready so we are waiting for this player
				g.Player(idx).AllowActions([]string{
					"pay",
				})
			}
		})

		// Reset player action states when task completed
		wp.OnCompleted(func() {
			g.ResetAllPlayerAllowedActions()
		})

		event.Payload.Task.AddTask(wp)

		return wp
	}

	return t
}

func (g *game) WaitForPayment(taskName string) bool {

	// Getting current event
	event := g.GetEvent()

	t := g.AssertPaymentTask(taskName)

	// Getting current task to execute if it is what we need
	if event.Payload.Task.GetAvailableTask() != t {
		return false
	}

	if !t.Execute() {
		//fmt.Println("Waiting for ante")
		return false
	}

	return true
}
