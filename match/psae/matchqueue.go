package psae

type MatchQueue interface {
	Publish(*Match) error
	Subscribe() (chan *Match, error)
	Close() error
}
