package slices

import (
	"fmt"
)

func Pop[T any](s *[]T) (T, error) {
	popped, err := PopN(s, 1)
	if err != nil {
		var empty T
		return empty, err
	}
	return popped[0], err
}

func PopN[T any](s *[]T, n int) ([]T, error) {
	size := len(*s)
	if size < n {
		var empty []T
		return empty, fmt.Errorf("Slice doesn't have enough elements to pop %d", n)
	}
	val := (*s)[size-n : size]
	(*s) = (*s)[:size-(n)]
	return val, nil
}

func Shift[T any](s *[]T) (T, error) {
	size := len(*s)
	if size == 0 {
		var empty T
		return empty, fmt.Errorf("Slice is empty")
	}
	val := (*s)[0]
	(*s) = (*s)[1:]
	return val, nil
}
