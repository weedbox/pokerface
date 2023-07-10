package psae

type WaitingRoom interface {
	Enter(PSAE, *Player) error
	Leave(PSAE, string) error
	Drain(PSAE) error
	Flush(PSAE) error
}
