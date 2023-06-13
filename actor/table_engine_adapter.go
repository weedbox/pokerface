package actor

import (
	"encoding/json"

	"github.com/weedbox/pokertable"
)

type tableEngineAdapter struct {
	actor        Actor
	engine       pokertable.TableEngine
	tableSetting pokertable.TableSetting
	table        *pokertable.Table
}

func NewTableEngineAdapter(te pokertable.TableEngine, table *pokertable.Table) *tableEngineAdapter {

	return &tableEngineAdapter{
		engine: te,
		table:  table,
	}
}

func (tea *tableEngineAdapter) SetActor(a Actor) {
	tea.actor = a
}

func (tea *tableEngineAdapter) UpdateTableState(tableInfo *pokertable.Table) error {

	// Clone to get a individual table structure
	data, err := tableInfo.GetJSON()
	if err != nil {
		return err
	}

	var t pokertable.Table
	err = json.Unmarshal([]byte(*data), &t)
	if err != nil {
		return err
	}

	tea.table = &t

	return tea.actor.UpdateTableState(&t)
}

func (tea *tableEngineAdapter) Pass(playerID string) error {
	return tea.engine.PlayerPass(tea.table.ID, playerID)
}

func (tea *tableEngineAdapter) Ready(playerID string) error {
	return tea.engine.PlayerReady(tea.table.ID, playerID)
}

func (tea *tableEngineAdapter) Pay(playerID string, chips int64) error {
	return tea.engine.PlayerPay(tea.table.ID, playerID, chips)
}

func (tea *tableEngineAdapter) Check(playerID string) error {
	return tea.engine.PlayerCheck(tea.table.ID, playerID)
}

func (tea *tableEngineAdapter) Bet(playerID string, chips int64) error {
	return tea.engine.PlayerBet(tea.table.ID, playerID, chips)
}

func (tea *tableEngineAdapter) Call(playerID string) error {
	return tea.engine.PlayerCall(tea.table.ID, playerID)
}

func (tea *tableEngineAdapter) Fold(playerID string) error {
	return tea.engine.PlayerFold(tea.table.ID, playerID)
}

func (tea *tableEngineAdapter) Allin(playerID string) error {
	return tea.engine.PlayerAllin(tea.table.ID, playerID)
}

func (tea *tableEngineAdapter) Raise(playerID string, chipLevel int64) error {
	return tea.engine.PlayerRaise(tea.table.ID, playerID, chipLevel)
}
