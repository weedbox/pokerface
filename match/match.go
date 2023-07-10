package match

import (
	"time"

	"github.com/weedbox/pokerface/match/psae"
)

type Match interface {
	DisallowRegistration()
	Join(playerID string) error
	Close() error
}

type MatchOpt func(*match)
type match struct {
	engine psae.PSAE
	b      Backend
}

func WithBackend(b Backend) MatchOpt {
	return func(m *match) {
		m.b = b
	}
}

func NewMatch(options *Options, opts ...MatchOpt) *match {

	m := &match{}

	for _, opt := range opts {
		opt(m)
	}

	// Preparing game configuration
	g := psae.NewGame()
	g.MaxPlayersPerTable = options.MaxSeats
	g.MinInitialPlayers = options.MinInitialPlayers
	g.TableLimit = options.MaxTables

	if !options.Joinable {
		g.Status = psae.GameStatus_AfterRegistrationDeadline
	}

	// Create backend
	b := m.initializeBackend()

	// Initializing allocation engine
	m.engine = psae.NewPSAE(
		psae.WithBackend(b),
		psae.WithRuntime(NewMatchRuntime()),
		psae.WithWaitingRoom(psae.NewMemoryWaitingRoom(time.Second*time.Duration(options.WaitingPeriod))),
		psae.WithGame(g),
		//TODO: implement SeatMap with Redis
		//TODO: implement WaitingRoom with Redis
		//TODO: implement MatchQueue with JetStream
		//TODO: implement DispatchQueue with JetStream
		//TODO: implement ReleaseQueue with JetStream
	)

	return m
}

func (m *match) initializeBackend() *psae.Backend {

	b := psae.NewBackend()

	b.AllocateTable = func() (*psae.TableState, error) {

		tableID, err := m.b.AllocateTable()
		if err != nil {
			return nil, err
		}

		ts := &psae.TableState{
			ID:             tableID,
			Players:        make(map[string]*psae.Player),
			Status:         psae.TableStatus_Ready,
			TotalSeats:     m.engine.Game().MaxPlayersPerTable,
			AvailableSeats: m.engine.Game().MaxPlayersPerTable,
			Statistics: &psae.TableStatistics{
				NoChanges: 0,
			},
		}

		return ts, nil
	}

	b.JoinTable = func(tableID string, players []*psae.Player) error {
		_, err := m.b.Join(tableID, players)
		return err
	}

	b.BrokeTable = func(tableID string) error {
		return m.b.BreakTable(tableID)
	}

	return b
}

func (m *match) Join(playerID string) error {

	player := &psae.Player{
		ID: playerID,
	}

	err := m.engine.Join(player)
	if err != nil {
		return err
	}

	return nil
}

func (m *match) Close() error {
	return m.engine.Close()
}

func (m *match) DisallowRegistration() {
	m.engine.DisallowRegistration()
}

func (m *match) GetTableState(tableId string) (*psae.TableState, error) {
	return m.engine.GetTableState(tableId)
}

func (m *match) UpdateTable(state *psae.TableState) (*psae.TableState, error) {
	return m.engine.UpdateTableState(state)
}
