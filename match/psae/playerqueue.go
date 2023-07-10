package psae

type PlayerQueue interface {
	Publish(*Player) error
	Subscribe() (chan *Player, error)
	Close() error
}
