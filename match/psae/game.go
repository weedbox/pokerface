package psae

type GameStatus int

const (
	GameStatus_Normal GameStatus = iota
	GameStatus_AfterRegistrationDeadline
	GameStatus_Suspend
)

type Game struct {
	Status             GameStatus
	MinInitialPlayers  int
	MaxPlayersPerTable int
	TableLimit         int
}

func NewGame() *Game {
	return &Game{
		Status:             GameStatus_Normal,
		MinInitialPlayers:  4,
		MaxPlayersPerTable: 9,
		TableLimit:         -1,
	}
}
