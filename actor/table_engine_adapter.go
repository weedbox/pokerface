package actor

import (
	pokertable "github.com/weedbox/pokertable"
	pokertableModel "github.com/weedbox/pokertable/model"
)

type tableEngineAdapter struct {
	actor        Actor
	engine       pokertable.TableEngine
	tableSetting pokertableModel.TableSetting
	table        *pokertableModel.Table
}

func NewTableEngineAdapter(te pokertable.TableEngine, table pokertableModel.Table) *tableEngineAdapter {

	return &tableEngineAdapter{
		engine: te,
		table:  &table,
	}
}

func (tea *tableEngineAdapter) SetActor(a Actor) {
	tea.actor = a
}

func (tea *tableEngineAdapter) UpdateTableState(tableInfo *pokertableModel.Table) error {
	return tea.actor.UpdateTableState(tableInfo)
}

func (tea *tableEngineAdapter) Ready(playerID string) error {

	_, err := tea.engine.PlayerReady(*tea.table, playerID)
	if err != nil {
		return err
	}

	return nil
}

func (tea *tableEngineAdapter) Pay(playerID string, chips int64) error {
	return tea.engine.PlayerPay(tea.table, playerID, chips)
}

func (tea *tableEngineAdapter) Check(playerID string) error {
	return tea.engine.PlayerCheck(tea.table, playerID)
}

func (tea *tableEngineAdapter) Bet(playerID string, chips int64) error {
	return tea.engine.PlayerBet(tea.table, playerID, chips)
}

func (tea *tableEngineAdapter) Call(playerID string) error {
	return tea.engine.PlayerCall(tea.table, playerID)
}

func (tea *tableEngineAdapter) Fold(playerID string) error {
	return tea.engine.PlayerFold(tea.table, playerID)
}

func (tea *tableEngineAdapter) Allin(playerID string) error {
	return tea.engine.PlayerAllin(tea.table, playerID)
}

func (tea *tableEngineAdapter) Raise(playerID string, chipLevel int64) error {
	return tea.engine.PlayerRaise(tea.table, playerID)
}
