package stack

import (
	"fmt"
)

type Stack[T any] []T

func (s *Stack[T]) Push(val ...T) {
	*s = append(*s, val...)
}

func (s *Stack[T]) Empty() bool {
	return s.Len() == 0
}

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) Pop() (T, error) {
	popped, err := s.PopMulti(1)
	if err != nil {
		var empty T
		return empty, err
	}
	return popped[0], err
}

func (s *Stack[T]) PopMulti(n int) ([]T, error) {
	if s.Len() < n {
		var empty []T
		return empty, fmt.Errorf("Stack doesn't have enough elements to pop %d", n)
	}
	val := (*s)[len(*s)-n : len(*s)]
	*s = (*s)[:len(*s)-(n)]
	return val, nil
}
