package table

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/weedbox/pokerface/seat_manager"
	"github.com/weedbox/syncsaga"
	"github.com/weedbox/timebank"
)

var (
	ErrRunningAlready              = errors.New("table: running already")
	ErrNotJoinable                 = errors.New("table: table is not joinable")
	ErrNotFoundPlayer              = errors.New("table: not found player")
	ErrPlayerNotInGame             = errors.New("table: player not in the game")
	ErrTimesUp                     = errors.New("table: time's up")
	ErrInsufficientNumberOfPlayers = errors.New("table: insufficient number of players")
	ErrGameConditionsNotMet        = errors.New("table: game conditions not met")
	ErrMaxGamesExceeded            = errors.New("table: reach the maximum number of games")
	ErrGameCancelled               = errors.New("table: game was cancelled")
)

type TableOpt func(*table)

type Table interface {
	Start() error
	Close() error
	Resume() error
	Pause() error

	// Player management
	Join(seatID int, p *PlayerInfo) (int, error)
	Leave(seatID int) error
	Reserve(seatID int) error
	Activate(seatID int) error
	ActivateByPlayerID(playerID string) error

	// Getter
	GetState() *State
	GetGame() Game
	GetGameCount() int
	GetPlayablePlayerCount() int
	GetPlayerByID(playerID string) *PlayerInfo
	GetPlayerByGameIdx(idx int) *PlayerInfo
	GetPlayerIdx(playerID string) int

	// Setter
	SetAnte(chips int64)
	SetBlinds(dealer int64, sb int64, bb int64)
	SetJoinable(enabled bool)

	// Event
	OnStateUpdated(func(*State))

	// Actions
	Ready(playerID string) error
	Pass(playerID string) error
	Pay(playerID string, chips int64) error
	Fold(playerID string) error
	Check(playerID string) error
	Call(playerID string) error
	Allin(playerID string) error
	Bet(playerID string, chips int64) error
	Raise(playerID string, chipLevel int64) error
}

type table struct {
	g              Game
	b              Backend
	isRunning      bool
	isPaused       bool
	inPosition     bool
	options        *Options
	gameCount      int
	gameLoop       chan int
	ts             *State
	rg             *syncsaga.ReadyGroup
	sm             *seat_manager.SeatManager
	tb             *timebank.TimeBank
	onStateUpdated func(*State)
}

func WithBackend(b Backend) TableOpt {
	return func(t *table) {
		t.b = b
	}
}

func NewTable(options *Options, opts ...TableOpt) *table {

	t := &table{
		options:        options,
		rg:             syncsaga.NewReadyGroup(),
		sm:             seat_manager.NewSeatManager(options.MaxSeats),
		ts:             NewState(),
		tb:             timebank.NewTimeBank(),
		gameLoop:       make(chan int, 1024),
		onStateUpdated: func(*State) {},
	}

	for _, opt := range opts {
		opt(t)
	}

	t.ts.ID = uuid.New().String()
	t.ts.Options = options
	t.ts.Status = "idle"

	return t
}

func (t *table) OnStateUpdated(fn func(*State)) {
	t.onStateUpdated = fn
}

func (t *table) GetState() *State {
	return t.ts
}

func (t *table) GetGame() Game {
	return t.g
}

func (t *table) GetGameCount() int {
	return t.gameCount
}

func (t *table) SetAnte(chips int64) {
	t.options.Ante = chips
}

func (t *table) SetBlinds(dealer int64, sb int64, bb int64) {
	t.options.Blind.Dealer = dealer
	t.options.Blind.SB = sb
	t.options.Blind.BB = bb
}

func (t *table) SetJoinable(enabled bool) {
	t.options.Joinable = enabled
}

func (t *table) Start() error {

	if t.isRunning {
		return nil
	}

	t.isRunning = true
	t.ts.StartTime = time.Now().Unix()
	t.ts.EndTime = t.ts.StartTime + int64(t.options.Duration)

	go t.tableLoop()

	//go t.nextGame(0)
	t.NewGame(0)

	return nil
}

func (t *table) NewGame(interval int) error {
	t.gameLoop <- interval
	return nil
}

func (t *table) Close() error {

	t.ts.Status = "closed"

	return nil
}

func (t *table) Resume() error {

	if !t.isPaused {
		return nil
	}

	t.isPaused = false
	t.ts.Status = "idle"

	if t.isRunning {
		return t.NewGame(0)
	}

	return nil
}

func (t *table) Pause() error {

	if t.isPaused {
		return nil
	}

	t.isPaused = true
	t.ts.Status = "pause"

	t.tb.Cancel()

	return nil
}

func (t *table) Activate(seatID int) error {

	err := t.sm.Activate(seatID)
	if err != nil {
		return nil
	}

	if !t.isRunning || t.ts.Status != "idle" {
		return nil
	}

	if t.sm.GetPlayableSeatCount() >= t.options.InitialPlayers {
		// Strarting game right now
		t.NewGame(0)
	}

	return nil
}

func (t *table) Reserve(seatID int) error {
	return t.sm.Reserve(seatID)
}

func (t *table) ActivateByPlayerID(playerID string) error {

	// Find the player before joining
	var found *PlayerInfo
	for _, ps := range t.ts.Players {
		if ps.ID == playerID {
			found = ps
		}
	}

	// Player is getting back to seat
	if found == nil {
		return ErrNotFoundPlayer
	}

	// Activate the seat
	err := t.Activate(found.SeatID)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Join(seatID int, p *PlayerInfo) (int, error) {

	sid, err := t.sm.Join(seatID, p)
	if err != nil {
		return -1, err
	}

	p.SeatID = sid
	t.ts.Players[sid] = p

	// Event
	state := t.cloneState()
	t.onStateUpdated(state)

	return sid, nil
}

func (t *table) Leave(seatID int) error {

	err := t.sm.Leave(seatID)
	if err != nil {
		return err
	}

	delete(t.ts.Players, seatID)

	// Event
	state := t.cloneState()
	t.onStateUpdated(state)

	return nil
}

func (t *table) GetPlayablePlayerCount() int {
	return t.sm.GetPlayableSeatCount()
}

func (t *table) GetPlayerByID(playerID string) *PlayerInfo {

	for _, p := range t.ts.Players {
		if p.ID == playerID {
			return p
		}
	}

	return nil
}

func (t *table) GetPlayerByGameIdx(idx int) *PlayerInfo {

	for _, p := range t.ts.Players {
		if p.GameIdx == idx {
			return p
		}
	}

	return nil
}

func (t *table) GetPlayerIdx(playerID string) int {

	p := t.GetPlayerByID(playerID)
	if p == nil {
		return -1
	}

	return p.GameIdx
}

// Actions
func (t *table) Ready(playerID string) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Ready(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Pass(playerID string) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Pass(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Pay(playerID string, chips int64) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Pay(idx, chips)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Fold(playerID string) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Fold(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Check(playerID string) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Check(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Call(playerID string) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Call(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Allin(playerID string) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Allin(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Bet(playerID string, chips int64) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Bet(idx, chips)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Raise(playerID string, chipLevel int64) error {

	if t.isPaused {
		return nil
	}

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Raise(idx, chipLevel)
	if err != nil {
		return err
	}

	return nil
}
