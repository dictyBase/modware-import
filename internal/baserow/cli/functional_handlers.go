package cli

import (
	"io/fs"
	"strings"

	F "github.com/IBM/fp-go/function"
)

var (
	Split = F.Curry2(F.Bind2nd[string, string, []string])(strings.Split)
)

type SliceWithError[T any] struct {
	Slice []T
	Error error
}

func onErrorWithSlice[T any](err error) SliceWithError[T] {
	return SliceWithError[T]{Error: err}
}

func onSuccessWithSlice[T any](slice []T) SliceWithError[T] {
	return SliceWithError[T]{Slice: slice}
}

func noDir(rec fs.DirEntry) bool { return !rec.IsDir() }
