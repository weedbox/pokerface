package psae

import (
	"math"
	"time"
)

type Opt func(*PSAEImpl)
type PSAEImpl struct {
	g       *Game
	backend *Backend
	sm      SeatMap
	wr      WaitingRoom
	rt      Runtime
	mq      MatchQueue
	dq      PlayerQueue
	rq      PlayerQueue
}

type PSAE interface {
	Game() *Game
	Backend() *Backend
	SeatMap() SeatMap
	WaitingRoom() WaitingRoom
	Runtime() Runtime
	MatchQueue() MatchQueue
	DispatchQueue() PlayerQueue
	ReleaseQueue() PlayerQueue
	Close() error

	// Game status
	GetStatus() GameStatus
	ResumeGame()
	SuspendGame()
	DisallowRegistration()
	IsLastTableStage() bool

	// Operations of table
	AllocateTable() (*TableState, error)
	BreakTable(string) error
	JoinTable(string, []*Player) error

	// Seatmap management
	AssertTableState(*TableState) error
	SetTableStatus(tid string, s TableStatus) error
	UpdateTableState(*TableState) (*TableState, error)
	GetTableState(id string) (*TableState, error)

	// Player arrangement
	Join(*Player) error
	DispatchPlayer(*Player) error
	ReleasePlayer(*Player) error
	MatchPlayers([]*Player) error

	// Waiting room
	EnterWaitingRoom(*Player) error
	LeaveWaitingRoom(string) error
	DrainWaitingRoom() error
	FlushWaitingRoom() error

	// Events
	EmitPlayerDispatched(*Player)
	EmitPlayerReleased(*Player)
	EmitWaitingRoomDrained(*Player)
	EmitWaitingRoomEntered(*Player)
	EmitWaitingRoomMatched([]*Player)
}

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewPSAE(opts ...Opt) PSAE {
	pi := &PSAEImpl{}

	for _, opt := range opts {
		opt(pi)
	}

	if pi.g == nil {
		pi.g = NewGame()
	}

	if pi.backend == nil {
		pi.backend = NewBackend()
	}

	if pi.sm == nil {
		pi.sm = NewMemorySeatMap()
	}

	if pi.wr == nil {
		pi.wr = NewMemoryWaitingRoom(time.Second * 15)
	}

	if pi.rt == nil {
		pi.rt = NewDefaultRuntime()
	}

	if pi.mq == nil {
		pi.mq = NewMemoryMatchQueue()
	}

	if pi.dq == nil {
		pi.dq = NewMemoryPlayerQueue()
	}

	if pi.rq == nil {
		pi.rq = NewMemoryPlayerQueue()
	}

	pi.init()

	return pi
}

func WithGame(g *Game) Opt {
	return func(pi *PSAEImpl) {
		pi.g = g
	}
}

func WithBackend(backend *Backend) Opt {
	return func(pi *PSAEImpl) {
		pi.backend = backend
	}
}

func WithSeatMap(sm SeatMap) Opt {
	return func(pi *PSAEImpl) {
		pi.sm = sm
	}
}

func WithWaitingRoom(wr WaitingRoom) Opt {
	return func(pi *PSAEImpl) {
		pi.wr = wr
	}
}

func WithRuntime(rt Runtime) Opt {
	return func(pi *PSAEImpl) {
		pi.rt = rt
	}
}

func WithDispatchQueue(dq PlayerQueue) Opt {
	return func(pi *PSAEImpl) {
		pi.dq = dq
	}
}

func WithReleaseQueue(rq PlayerQueue) Opt {
	return func(pi *PSAEImpl) {
		pi.rq = rq
	}
}

func (pi *PSAEImpl) Game() *Game {
	return pi.g
}

func (pi *PSAEImpl) Backend() *Backend {
	return pi.backend
}

func (pi *PSAEImpl) SeatMap() SeatMap {
	return pi.sm
}

func (pi *PSAEImpl) init() error {

	// Initialize match queue
	mq, err := pi.MatchQueue().Subscribe()
	if err != nil {
		return err
	}

	go func() {
		// Waiting for matching
		for m := range mq {
			pi.EmitMatched(m)
		}
	}()

	// Initialize player dispatcher
	pdq, err := pi.DispatchQueue().Subscribe()
	if err != nil {
		return err
	}

	go func() {
		// Waiting for dispatched players
		for player := range pdq {
			pi.EmitPlayerDispatched(player)
		}
	}()

	// Initialize queue for released players
	prq, err := pi.ReleaseQueue().Subscribe()
	if err != nil {
		return err
	}

	go func() {
		// Waiting for released players
		for player := range prq {
			pi.EmitPlayerReleased(player)
		}
	}()

	return nil
}

func (pi *PSAEImpl) WaitingRoom() WaitingRoom {
	return pi.wr
}

func (pi *PSAEImpl) Runtime() Runtime {
	return pi.rt
}

func (pi *PSAEImpl) MatchQueue() MatchQueue {
	return pi.mq
}

func (pi *PSAEImpl) DispatchQueue() PlayerQueue {
	return pi.dq
}

func (pi *PSAEImpl) ReleaseQueue() PlayerQueue {
	return pi.rq
}

func (pi *PSAEImpl) Close() error {

	err := pi.ReleaseQueue().Close()
	if err != nil {
		return err
	}

	err = pi.DispatchQueue().Close()
	if err != nil {
		return err
	}

	err = pi.MatchQueue().Close()
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) GetStatus() GameStatus {
	return pi.g.Status
}

func (pi *PSAEImpl) IsLastTableStage() bool {

	if pi.GetStatus() == GameStatus_AfterRegistrationDeadline {

		count, err := pi.SeatMap().GetTableCount()
		if err != nil {
			return false
		}

		// Final table
		if count == 1 {
			return true
		}

		totalPlayers, err := pi.SeatMap().GetTotalPlayers()
		if err != nil {
			return false
		}

		// too many tables
		proposedTableCount := math.Ceil(float64(totalPlayers) / float64(pi.Game().MaxPlayersPerTable))
		if proposedTableCount == 1 {
			return true
		}
	}

	return false
}

func (pi *PSAEImpl) ResumeGame() {
	pi.g.Status = GameStatus_Normal
}

func (pi *PSAEImpl) SuspendGame() {
	pi.g.Status = GameStatus_Suspend
}

func (pi *PSAEImpl) DisallowRegistration() {
	pi.g.Status = GameStatus_AfterRegistrationDeadline
}

func (pi *PSAEImpl) Join(player *Player) error {
	return pi.DispatchPlayer(player)
}

func (pi *PSAEImpl) AllocateTable() (*TableState, error) {
	return pi.Backend().AllocateTable()
}

func (pi *PSAEImpl) AssertTableState(ts *TableState) error {

	origts, err := pi.sm.GetTableState(ts.ID)
	if err != nil {
		return err
	}

	// Not exist
	if origts == nil {
		err := pi.sm.CreateTable(ts)
		if err != nil {
			return err
		}
	}

	// Update table
	_, err = pi.sm.UpdateTable(ts)
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) GetTableState(tId string) (*TableState, error) {

	ts, err := pi.sm.GetTableState(tId)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (pi *PSAEImpl) SetTableStatus(tid string, s TableStatus) error {

	ts, err := pi.sm.GetTableState(tid)
	if err != nil {
		return err
	}

	ts.Status = s

	_, err = pi.sm.UpdateTable(ts)
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) UpdateTableState(ts *TableState) (*TableState, error) {

	// Switch to be suspended
	ts.Status = TableStatus_Busy

	// Update table
	newts, err := pi.sm.UpdateTable(ts)
	if err != nil {
		return nil, err
	}

	// Emit event to trigger runtime to check state
	pi.rt.TableStateUpdated(pi, newts)

	if newts.Status != TableStatus_Busy {
		// Table is broken
		return newts, nil
	}

	// Back to ready
	newts.Status = TableStatus_Ready

	// Update table
	newts, err = pi.sm.UpdateTable(newts)
	if err != nil {
		return nil, err
	}

	return newts, nil
}

func (pi *PSAEImpl) BreakTable(tid string) error {

	ts, err := pi.sm.GetTableState(tid)
	if err != nil {
		return err
	}

	// call API to broke table
	err = pi.Backend().BrokeTable(tid)
	if err != nil {
		return err
	}

	err = pi.sm.DestroyTable(tid)
	if err != nil {
		return err
	}

	// emit event
	pi.rt.TableBroken(pi, ts)

	return nil
}

func (pi *PSAEImpl) JoinTable(tid string, players []*Player) error {

	if len(players) == 0 {
		return nil
	}

	err := pi.Backend().JoinTable(tid, players)
	if err != nil {
		return err
	}

	ts, err := pi.sm.GetTableState(tid)
	if err != nil {
		return err
	}

	// Update player list
	for _, p := range players {
		ts.Players[p.ID] = p
	}

	//TODO: update by following rules
	ts.AvailableSeats = ts.TotalSeats - len(ts.Players)

	// Update table state
	_, err = pi.sm.UpdateTable(ts)
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) DispatchPlayer(p *Player) error {
	return pi.dq.Publish(p)
}

func (pi *PSAEImpl) ReleasePlayer(p *Player) error {
	//fmt.Printf("Release player: %s\n", p.ID)
	return pi.rq.Publish(p)
}

func (pi *PSAEImpl) MatchPlayers(players []*Player) error {

	m := &Match{
		Players: make([]*Player, len(players)),
	}

	for i, p := range players {
		m.Players[i] = p
	}

	return pi.mq.Publish(m)
}

func (pi *PSAEImpl) EnterWaitingRoom(player *Player) error {

	err := pi.WaitingRoom().Enter(pi, player)
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) LeaveWaitingRoom(pid string) error {

	err := pi.WaitingRoom().Leave(pi, pid)
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) DrainWaitingRoom() error {

	err := pi.WaitingRoom().Drain(pi)
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) FlushWaitingRoom() error {

	err := pi.WaitingRoom().Flush(pi)
	if err != nil {
		return err
	}

	return nil
}

func (pi *PSAEImpl) EmitMatched(m *Match) {
	pi.rt.Matched(pi, m)
}

func (pi *PSAEImpl) EmitPlayerDispatched(player *Player) {
	pi.rt.PlayerDispatched(pi, player)
}

func (pi *PSAEImpl) EmitPlayerReleased(player *Player) {
	pi.rt.PlayerReleased(pi, player)
}

func (pi *PSAEImpl) EmitWaitingRoomDrained(player *Player) {
	pi.rt.WaitingRoomDrained(pi, player)
}

func (pi *PSAEImpl) EmitWaitingRoomEntered(player *Player) {
	pi.rt.WaitingRoomEntered(pi, player)
}

func (pi *PSAEImpl) EmitWaitingRoomLeft(player *Player) {
	pi.rt.WaitingRoomLeft(pi, player)
}

func (pi *PSAEImpl) EmitWaitingRoomMatched(players []*Player) {
	pi.rt.WaitingRoomMatched(pi, players)
}
