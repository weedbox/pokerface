package psae

type SeatMap interface {
	CreateTable(tableState *TableState) error
	DestroyTable(tableID string) error
	UpdateTable(tableState *TableState) (*TableState, error)
	GetTableState(tableID string) (*TableState, error)
	FindAvailableTable(condition *TableCondition) (*TableState, error)
	GetTableCount() (int, error)
	GetTotalPlayers() (int64, error)
	GetAllTables() ([]*TableState, error)
}

type TableStatus int

const (
	TableStatus_Ready TableStatus = iota
	TableStatus_Busy
	TableStatus_Broken
	TableStatus_Suspend
)

type TableState struct {
	ID             string             `json:"id"`
	Players        map[string]*Player `json:"players"`
	Status         TableStatus        `json:"status"`
	AvailableSeats int                `json:"availableSeats"`
	TotalSeats     int                `json:"totalSeats"`
	LastGameID     string             `json:"lastGameID"`
	Statistics     *TableStatistics   `json:"statistics",omitempty`
}

type TableStatistics struct {
	NoChanges int
}

type TableCondition struct {
	HighestNumberOfPlayers bool
	MinAvailableSeats      int
}

func (ts *TableState) Clone() *TableState {

	newts := &TableState{
		ID:             ts.ID,
		Players:        make(map[string]*Player),
		Status:         ts.Status,
		AvailableSeats: ts.AvailableSeats,
		TotalSeats:     ts.TotalSeats,
		LastGameID:     ts.LastGameID,
		Statistics: &TableStatistics{
			NoChanges: ts.Statistics.NoChanges,
		},
	}

	for pid, p := range ts.Players {
		newts.Players[pid] = p
	}

	return newts
}
