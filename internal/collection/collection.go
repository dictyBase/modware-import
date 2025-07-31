package collection

import (
	"slices"
)

// Filter returns a new slice containing all elements that satisfy the
// predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

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

// Pipe5 creates a functional pipeline by taking an initial value and applying
// five functions in succession. The output of each function becomes the input
// to the next function. The final return value is the result of the last
// function application.
func Pipe5[T1, T2, T3, T4, T5, T6 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
) T6 {
	return fn5(fn4(fn3(fn2(fn1(initial)))))
}

// Pipe4 creates a functional pipeline by taking an initial value and applying
// four functions in succession. The output of each function becomes the input
// to the next function. The final return value is the result of the last
// function application.
func Pipe4[T1, T2, T3, T4, T5 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
) T5 {
	return fn4(fn3(fn2(fn1(initial))))
}

// Pipe2 creates a functional pipeline by taking an initial value and applying
// two functions in succession. The output of the first function becomes the
// input to the second function. The final return value is the result of the
// last function application.
func Pipe2[T1, T2, T3 any](tup T1, f1 func(T1) T2, fn2 func(T2) T3) T3 {
	return fn2(f1(tup))
}

// CurriedMap returns a function that, when given a slice, applies the provided
// function to each element of the slice. This is a curried version of the Map
// function.
func CurriedMap[T1, T2 any](fnc func(T1) T2) func([]T1) []T2 {
	return func(slc []T1) []T2 {
		return Map(slc, fnc)
	}
}

// CurriedFilter returns a function that, when given a slice, filters it based on
// the provided predicate. This is a curried version of the Filter function.
func CurriedFilter[T any](predicate func(T) bool) func([]T) []T {
	return func(slice []T) []T {
		return Filter(slice, predicate)
	}
}
