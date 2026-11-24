//go:build !linux && !darwin

package normalize

import "errors"

var errNoAtomicRename = errors.New(
	"atomic no-replace directory rename unavailable on this platform",
)

// renameNew fails, because this platform has no no-replace rename. Do not
// use os.Lstat followed by os.Rename: two runs can race.
func renameNew(string, string) (string, error) { return "", errNoAtomicRename }
