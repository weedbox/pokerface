package actor

import (
	"math/rand"
	"time"

	"github.com/weedbox/pokerface"
	"github.com/weedbox/pokertable"
)

type ActionProbability struct {
	Action string
	Weight float64
}

var (
	actionProbabilities = []ActionProbability{
		{Action: "check", Weight: 0.1},
		{Action: "call", Weight: 0.3},
		{Action: "fold", Weight: 0.2},
		{Action: "allin", Weight: 0.1},
		{Action: "raise", Weight: 0.2},
		{Action: "bet", Weight: 0.1},
	}
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
	br.gamePlayerIdx = table.GamePlayerIndex(br.playerID)
	br.tableInfo = table

	// Somehow, this player is not in the game.
	// It probably has no chips already.
	if br.gamePlayerIdx == -1 {
		return nil
	}

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

	//fmt.Println(br.tableInfo.State.GameState.Status.Round, br.gamePlayerIdx, gs.Players[br.gamePlayerIdx].AllowedActions)

	// Do ready() and pay() automatically
	if gs.HasAction(br.gamePlayerIdx, "ready") {
		return br.actions.Ready()
	} else if gs.HasAction(br.gamePlayerIdx, "pass") {
		return br.actions.Pass()
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

func (br *botRunner) calcActionProbabilities(actions []string) map[string]float64 {

	probabilities := make(map[string]float64)
	totalWeight := 0.0
	for _, action := range actions {

		for _, p := range actionProbabilities {
			if action == p.Action {
				probabilities[action] = p.Weight
				totalWeight += p.Weight
				break
			}
		}
	}

	scaleRatio := 1.0 / totalWeight
	weightLevel := 0.0
	for action, weight := range probabilities {
		scaledWeight := weight * scaleRatio
		weightLevel += scaledWeight
		probabilities[action] = weightLevel
	}

	return probabilities
}

func (br *botRunner) calcAction(actions []string) string {

	// Select action randomly
	rand.Seed(time.Now().UnixNano())

	probabilities := br.calcActionProbabilities(actions)
	randomNum := rand.Float64()

	for action, probability := range probabilities {
		if randomNum < probability {
			return action
		}
	}

	return actions[0]
}

func (br *botRunner) requestAI() error {

	gs := br.tableInfo.State.GameState
	player := gs.Players[br.gamePlayerIdx]

	//TODO: To simulate human-like behavior, it is necessary to incorporate random delays when performing actions.

	// None of actions is allowed
	if len(player.AllowedActions) == 0 {
		return nil
	}

	//	actionIdx := 0
	action := player.AllowedActions[0]

	if len(player.AllowedActions) > 1 {
		action = br.calcAction(player.AllowedActions)
	}

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
