package actor

import (
	"math/rand"

	"github.com/weedbox/pokerface"
	"github.com/weedbox/pokertable"
)

type botRunner struct {
	actor         Actor
	actions       Actions
	playerID      string
	gamePlayerIdx int
	tableInfo     *pokertable.Table
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

	// Update player index in game
	br.gamePlayerIdx = table.PlayingPlayerIndex(br.playerID)
	br.tableInfo = table

	// Game is running right now
	switch br.tableInfo.State.Status {
	case pokertable.TableStateStatus_TableGameMatchOpen:

		// We have actions allowed by game engine
		player := br.tableInfo.State.GameState.GetPlayer(br.gamePlayerIdx)
		if len(player.AllowedActions) > 0 {
			br.requestMove()
		}
	}

	return nil
}

func (br *botRunner) requestMove() error {

	gs := br.tableInfo.State.GameState

	// Do ready() and pay() automatically
	if gs.HasAction(br.gamePlayerIdx, "ready") {
		return br.actions.Ready()
	} else if gs.HasAction(br.gamePlayerIdx, "pay") {

		// Pay for ante and blinds
		switch gs.Status.CurrentEvent.Name {
		case pokerface.GameEventSymbols[pokerface.GameEvent_Prepared]:

			// Ante
			return br.actions.Pay(gs.Meta.Ante)

		case pokerface.GameEventSymbols[pokerface.GameEvent_RoundInitialized]:

			// blinds
			if gs.HasPosition(br.gamePlayerIdx, "sb") {
				return br.actions.Pay(gs.Meta.Blind.SB)
			} else if gs.HasPosition(br.gamePlayerIdx, "bb") {
				return br.actions.Pay(gs.Meta.Blind.BB)
			}

			return br.actions.Pay(gs.Meta.Blind.Dealer)
		}
	}

	return br.requestAI()
}

func (br *botRunner) requestAI() error {

	gs := br.tableInfo.State.GameState
	player := gs.Players[br.gamePlayerIdx]

	//TODO: To simulate human-like behavior, it is necessary to incorporate random delays when performing actions.

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
