package collection

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
func Remove(a []string, items ...string) []string {
	var s []string
	for _, v := range a {
		if !Contains(items, v) {
			s = append(s, v)
		}
	}
	return s
}

// Contains reports whether s is present in slice a
func Contains(a []string, s string) bool {
	if len(a) == 0 {
		return false
	}
	return Index(a, s) >= 0
}

// Index returns the index of the first instance of s in slice a, or -1 if s is
// not present in a
func Index(a []string, s string) int {
	if len(a) == 0 {
		return -1
	}
	for i, v := range a {
		if v == s {
			return i
		}
	}
	return -1
}

// Extend appends all the elements of slices to a new slice
func Extend[T any](elems ...[]T) []T {
	aslice := make([]T, 0)
	for _, parts := range elems {
		aslice = append(aslice, parts...)
	}
	return aslice
}
