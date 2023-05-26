package pokerface

import (
	"fmt"

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
		fmt.Println("Waiting for ready")
		return false
	}

	return true
}
