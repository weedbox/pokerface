package actor

import (
	pokertable "github.com/weedbox/pokertable/model"
)

type Runner interface {
	SetActor(a Actor)
	UpdateTableState(t *pokertable.Table) error
}
