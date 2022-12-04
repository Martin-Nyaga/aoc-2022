package util

import "errors"

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

func (s *Set[T]) Remove(b T) bool {
	if !s.Has(b) {
		return false
	}
	delete(*s, b)
	return true
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
