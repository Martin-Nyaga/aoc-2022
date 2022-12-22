package set

import (
	"errors"
	"fmt"
)

type void struct{}
type Set[T comparable] map[T]void

func NewSet[T comparable](bytes ...T) Set[T] {
	set := Set[T]{}
	for _, b := range bytes {
		set[b] = void{}
	}
	return set
}

func (s *Set[T]) Add(b T) {
	(*s)[b] = void{}
}

func (s *Set[T]) Has(b T) bool {
	_, exists := (*s)[b]
	return exists
}

func (s *Set[T]) Empty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Len() int {
	return len(*s)
}

func (s *Set[T]) Remove(b T) error {
	if !s.Has(b) {
		return fmt.Errorf("Set doesn't have %#v", b)
	}
	delete(*s, b)
	return nil
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	intersection := Set[T]{}
	for k := range *s {
		if other.Has(k) {
			intersection.Add(k)
		}
	}
	return &intersection
}

func (s *Set[T]) PopAny() (T, error) {
	for k := range *s {
		s.Remove(k)
		return k, nil
	}
	var empty T
	return empty, errors.New("Set didn't have any elements to pop")
}

func (s *Set[T]) Each(fn func(T)) {
	for k := range *s {
		fn(k)
	}
}

func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, s.Len())
	s.Each(func(el T) {
		result = append(result, el)
	})
	return result
}
