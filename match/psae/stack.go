package psae

import (
	"container/list"
	"sync"
)

type Stack struct {
	list *list.List
	mu   sync.RWMutex
}

func NewStack() *Stack {
	return &Stack{
		list: list.New(),
	}
}

func (stack *Stack) List() *list.List {
	return stack.list
}

func (stack *Stack) Push(value interface{}) {
	stack.mu.Lock()
	defer stack.mu.Unlock()
	stack.list.PushBack(value)
}

func (stack *Stack) Pop() interface{} {
	stack.mu.Lock()
	defer stack.mu.Unlock()
	e := stack.list.Back()
	if e != nil {
		stack.list.Remove(e)
		return e.Value
	}
	return nil
}

func (stack *Stack) Peak() interface{} {
	e := stack.list.Back()
	if e != nil {
		return e.Value
	}

	return nil
}

func (stack *Stack) Len() int {
	return stack.list.Len()
}

func (stack *Stack) Empty() bool {
	return stack.list.Len() == 0
}
