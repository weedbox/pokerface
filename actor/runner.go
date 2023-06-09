package actor

import (
	"github.com/weedbox/pokertable"
)

type Runner interface {
	SetActor(a Actor)
	UpdateTableState(t *pokertable.Table) error
}
