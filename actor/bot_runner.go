package actor

import (
	"math/rand"

	pokertable "github.com/weedbox/pokertable/model"
)

type botRunner struct {
	actor     Actor
	actions   Actions
	playerID  string
	seatIndex int
	tableInfo *pokertable.Table
}

func NewBotRunner(playerID string) *botRunner {
	return &botRunner{
		playerID: playerID,
	}
}

func (br *botRunner) SetActor(a Actor) {
	br.actor = a
	br.actions = NewActions(a, br.playerID)
}

func (br *botRunner) UpdateTableState(table *pokertable.Table) error {

	// Update seat index
	br.seatIndex = br.actor.GetSeatIndex(br.playerID)
	br.tableInfo = table

	// Game is running right now
	switch br.tableInfo.State.Status {
	case pokertable.TableStateStatus_TableGameMatchOpen:

		// We have actions allowed by game engine
		player := br.tableInfo.State.GameState.GetPlayer(br.seatIndex)
		if len(player.AllowedActions) > 0 {
			br.requestMove()
		}
	}

	return nil
}

func (br *botRunner) requestMove() error {

	//TODO: To simulate human-like behavior, it is necessary to incorporate random delays when performing actions.

	gs := br.tableInfo.State.GameState
	player := gs.Players[br.seatIndex]

	//TODO: require smarter AI
	// Select action randomly
	actionIdx := rand.Intn(len(player.AllowedActions) - 1)
	action := player.AllowedActions[actionIdx]

	// Calculate chips
	switch action {
	case "bet":

		minBet := gs.Status.MiniBet
		chips := rand.Int63n(player.InitialStackSize-minBet) + minBet

		return br.actions.Bet(chips)
	case "raise":

		maxChipLevel := player.InitialStackSize
		minChipLevel := gs.Status.CurrentWager + gs.Status.PreviousRaiseSize

		chips := rand.Int63n(maxChipLevel-minChipLevel) + minChipLevel

		return br.actions.Raise(chips)
	case "call":
		return br.actions.Call()
	case "check":
		return br.actions.Check()
	case "allin":
		return br.actions.Allin()
	}

	return br.actions.Fold()
}
