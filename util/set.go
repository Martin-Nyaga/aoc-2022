package util

import "errors"

type void struct{}
type ByteSet map[byte]void

func NewByteSet(bytes ...byte) ByteSet {
	set := ByteSet{}
	for _, b := range bytes {
		set[b] = void{}
	}
	return set
}

func (s *ByteSet) Add(b byte) {
	(*s)[b] = void{}
}

func (s *ByteSet) Has(b byte) bool {
	_, exists := (*s)[b]
	return exists
}

func (s *ByteSet) Remove(b byte) bool {
	if !s.Has(b) {
		return false
	}
	delete(*s, b)
	return true
}

func (s *ByteSet) Intersection(other *ByteSet) *ByteSet {
	intersection := ByteSet{}
	for k := range *s {
		if other.Has(k) {
			intersection.Add(k)
		}
	}
	return &intersection
}

func (s *ByteSet) PopAny() (byte, error) {
	for k := range *s {
		s.Remove(k)
		return k, nil
	}
	return 0, errors.New("Set didn't have any elements to pop")
}
