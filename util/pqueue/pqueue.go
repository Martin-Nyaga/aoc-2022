package pqueue

import (
	"errors"

	"golang.org/x/exp/constraints"
)

type queuable interface {
	comparable
	constraints.Ordered
}

const (
	MinQueue = iota
	MaxQueue
)

type Pqueue[P queuable, V any] struct {
	elements   map[P][]V
	priorities []P
	mode       int
}

func NewPqueue[P queuable, V any](mode int) Pqueue[P, V] {
	return Pqueue[P, V]{
		elements:   map[P][]V{},
		priorities: make([]P, 0),
		mode:       mode,
	}
}

func (q *Pqueue[P, V]) Empty() bool {
	return q.Len() == 0
}

func (q *Pqueue[P, V]) Len() int {
	size := 0
	for _, arr := range q.elements {
		size += len(arr)
	}
	return size
}

func (q *Pqueue[P, V]) Push(p P, el V) {
	if q.hasPriority(p) {
		q.elements[p] = append(q.elements[p], el)
	} else {
		q.addPriority(p)
		q.elements[p] = []V{el}
	}
}

func (q *Pqueue[P, V]) Pop() (V, error) {
	var result V
	if q.Empty() {
		return result, errors.New("Queue was epmty")
	}
	priority := q.priorities[0]
	result = q.popValueWithPriority(priority)
	if _, exists := q.elements[priority]; !exists {
		q.priorities = q.priorities[1:]
	}
	return result, nil
}

func (q *Pqueue[P, V]) popValueWithPriority(p P) V {
	result := q.elements[p][0]
	if len(q.elements[p]) == 1 {
		delete(q.elements, p)
	} else {
		q.elements[p] = q.elements[p][1:]
	}
	return result
}

func (q *Pqueue[P, V]) hasPriority(priority P) bool {
	for _, prio := range q.priorities {
		if prio == priority {
			return true
		}
	}
	return false
}

func (q *Pqueue[P, V]) addPriority(priority P) {
	if len(q.priorities) == 0 {
		q.priorities = append(q.priorities, priority)
		return
	}

	insertIndex := 0
	for i, prio := range q.priorities {
		var insertHere bool
		if q.mode == MinQueue {
			insertHere = prio > priority
		} else {
			insertHere = prio < priority
		}
		if insertHere {
			insertIndex = i
			break
		}
		insertIndex += 1
	}

	newPriorities := make([]P, 0, len(q.priorities)+1)
	newPriorities = append(newPriorities, q.priorities[:insertIndex]...)
	newPriorities = append(newPriorities, priority)
	newPriorities = append(newPriorities, q.priorities[insertIndex:]...)
	q.priorities = newPriorities
}
