package collection

import (
	"golang.org/x/exp/slices"
)

// Map returns the slice obtained after applying the given function over every
// element in the given slice
func Map[T1, T2 any](s []T1, fn func(T1) T2) []T2 {
	ret := make([]T2, 0, len(s))
	for _, e := range s {
		ret = append(ret, fn(e))
	}
	return ret
}

// Remove removes items from the given(a) slice
func Remove[T comparable](a []T, items ...T) []T {
	var s []T
	for _, v := range a {
		if !slices.Contains(items, v) {
			s = append(s, v)
		}
	}
	return s
}

// Extend appends all the elements of slices to a new slice
func Extend[T any](elems ...[]T) []T {
	aslice := make([]T, 0)
	for _, parts := range elems {
		aslice = append(aslice, parts...)
	}
	return aslice
}
